package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func taskFetchWallets(client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("mangopay.taskFetchWallets")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
		)

		if err := fetchWallets(ctx, client, userID, connectorID, scheduler, ingester); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchWallets(
	ctx context.Context,
	client *client.Client,
	userID string,
	connectorID models.ConnectorID,
	scheduler task.Scheduler,
	ingester ingestion.Ingester,
) error {
	for page := 1; ; page++ {
		pagedWallets, err := client.GetWallets(ctx, userID, page)
		if err != nil {
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
					Reference:   wallet.ID,
					ConnectorID: connectorID,
				},
				CreatedAt:    time.Unix(wallet.CreationDate, 0),
				Reference:    wallet.ID,
				ConnectorID:  connectorID,
				DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, wallet.Currency),
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
					Reference:   wallet.ID,
					ConnectorID: connectorID,
				},
				Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, wallet.Balance.Currency),
				Balance:       &amount,
				CreatedAt:     now,
				LastUpdatedAt: now,
				ConnectorID:   connectorID,
			})
		}

		if err := ingester.IngestAccounts(ctx, accountBatch); err != nil {
			return err
		}

		if err := ingester.IngestBalances(ctx, balanceBatch, false); err != nil {
			return err
		}

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
