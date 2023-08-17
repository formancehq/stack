package modulr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/formancehq/payments/internal/app/connectors/modulr/client"
	"github.com/formancehq/payments/internal/app/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

var (
	accountsAndBalancesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts_and_balances"))...)
	accountsAttrs            = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))...)
	balancesAttrs            = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "balances"))...)
)

func taskFetchAccounts(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchAccounts)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), accountsAndBalancesAttrs)
		}()

		accounts, err := client.GetAccounts()
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAndBalancesAttrs)
			return err
		}

		if err := ingestAccountsBatch(ctx, ingester, metricsRegistry, accounts); err != nil {
			return err
		}

		for _, account := range accounts {
			logger.Infof("scheduling fetch-transactions: %s", account.ID)

			transactionsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch transactions from client by account",
				Key:       taskNameFetchTransactions,
				AccountID: account.ID,
			})
			if err != nil {
				return err
			}

			err = scheduler.Schedule(ctx, transactionsTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return err
			}
		}

		return nil
	}
}

func ingestAccountsBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	metricsRegistry metrics.MetricsRegistry,
	accounts []*client.Account,
) error {
	accountsBatch := ingestion.AccountBatch{}
	balancesBatch := ingestion.BalanceBatch{}

	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		openingDate, err := time.Parse("2006-01-02T15:04:05.999999999+0000", account.CreatedDate)
		if err != nil {
			return err
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference: account.ID,
				Provider:  models.ConnectorProviderModulr,
			},
			CreatedAt:    openingDate,
			Reference:    account.ID,
			Provider:     models.ConnectorProviderModulr,
			DefaultAsset: models.Asset(fmt.Sprintf("%s/2", account.Currency)),
			AccountName:  account.Name,
			Type:         models.AccountTypeInternal,
			RawData:      raw,
		})

		var amount big.Float
		_, ok := amount.SetString(account.Balance)
		if !ok {
			return fmt.Errorf("failed to parse amount %s", account.Balance)
		}

		var balance big.Int
		amount.Mul(&amount, big.NewFloat(100)).Int(&balance)

		now := time.Now()
		balancesBatch = append(balancesBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference: account.ID,
				Provider:  models.ConnectorProviderModulr,
			},
			Asset:         models.Asset(fmt.Sprintf("%s/2", account.Currency)),
			Balance:       &balance,
			CreatedAt:     now,
			LastUpdatedAt: now,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountsBatch)), accountsAttrs)

	if err := ingester.IngestBalances(ctx, balancesBatch, false); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(balancesBatch)), balancesAttrs)

	return nil
}
