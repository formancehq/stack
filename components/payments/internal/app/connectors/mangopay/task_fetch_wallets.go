package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/connectors/mangopay/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	walletsAndBalancesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "wallets_and_balances"))...)
	walletsAttrs            = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "wallets"))...)
	balancesAttrs           = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "balances"))...)
)

func taskFetchWallets(logger logging.Logger, client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchWallets)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), walletsAndBalancesAttrs)
		}()

		for page := 1; ; page++ {
			pagedWallets, err := client.GetWallets(ctx, userID, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, walletsAndBalancesAttrs)
				return err
			}

			if len(pagedWallets) == 0 {
				break
			}

			var accountBatch ingestion.AccountBatch
			var balanceBatch ingestion.BalanceBatch
			var transactionTasks []models.TaskDescriptor
			for _, wallet := range pagedWallets {
				transactionTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
					Name:     "Fetch transactions from client by user and wallets",
					Key:      taskNameFetchTransactions,
					UserID:   userID,
					WalletID: wallet.ID,
				})
				if err != nil {
					return err
				}

				buf, err := json.Marshal(wallet)
				if err != nil {
					return err
				}

				transactionTasks = append(transactionTasks, transactionTask)
				accountBatch = append(accountBatch, &models.Account{
					ID: models.AccountID{
						Reference: wallet.ID,
						Provider:  models.ConnectorProviderMangopay,
					},
					CreatedAt:    time.Unix(wallet.CreationDate, 0),
					Reference:    wallet.ID,
					Provider:     models.ConnectorProviderMangopay,
					DefaultAsset: currency.FormatAsset(wallet.Currency),
					AccountName:  wallet.Description,
					// Wallets are internal accounts on our side, since we
					// can have their balances.
					Type: models.AccountTypeInternal,
					Metadata: map[string]string{
						"user_id": userID,
					},
					RawData: buf,
				})

				var amount big.Int
				_, ok := amount.SetString(wallet.Balance.Amount.String(), 10)
				if !ok {
					return fmt.Errorf("failed to parse amount: %s", wallet.Balance.Amount.String())
				}

				now := time.Now()
				balanceBatch = append(balanceBatch, &models.Balance{
					AccountID: models.AccountID{
						Reference: wallet.ID,
						Provider:  models.ConnectorProviderMangopay,
					},
					Asset:         currency.FormatAsset(wallet.Balance.Currency),
					Balance:       &amount,
					CreatedAt:     now,
					LastUpdatedAt: now,
				})
			}

			if err := ingester.IngestAccounts(ctx, accountBatch); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, walletsAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountBatch)), walletsAttrs)

			if err := ingester.IngestBalances(ctx, balanceBatch, false); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(balanceBatch)), balancesAttrs)

			for _, transactionTask := range transactionTasks {
				err = scheduler.Schedule(ctx, transactionTask, models.TaskSchedulerOptions{
					ScheduleOption: models.OPTIONS_RUN_NOW,
					RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
				})
				if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
					return err
				}
			}
		}

		return nil
	}
}
