package currencycloud

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currencycloud/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	accountsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))...)
)

func taskFetchAccounts(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchAccounts)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), accountsAttrs)
		}()

		page := 1
		for {
			if page < 0 {
				break
			}

			pagedAccounts, nextPage, err := client.GetAccounts(ctx, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
				return err
			}

			page = nextPage

			if err := ingestAccountsBatch(ctx, ingester, pagedAccounts); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedAccounts)), accountsAttrs)
		}

		taskTransactions, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch transactions from client",
			Key:  taskNameFetchTransactions,
		})
		if err != nil {
			return err
		}

		err = scheduler.Schedule(ctx, taskTransactions, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}

		taskBalances, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch balances from client",
			Key:  taskNameFetchBalances,
		})
		if err != nil {
			return err
		}

		err = scheduler.Schedule(ctx, taskBalances, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}

		return nil
	}
}

func ingestAccountsBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	accounts []*client.Account,
) error {
	batch := ingestion.AccountBatch{}
	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		batch = append(batch, &models.Account{
			ID: models.AccountID{
				Reference: account.ID,
				Provider:  models.ConnectorProviderCurrencyCloud,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:   account.CreatedAt,
			Reference:   account.ID,
			Provider:    models.ConnectorProviderCurrencyCloud,
			AccountName: account.AccountName,
			Type:        models.AccountTypeInternal,
			RawData:     raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
