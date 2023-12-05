package moneycorp

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	accountsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))...)
)

func taskFetchAccounts(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchAccounts)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), accountsAttrs)
		}()

		for page := 1; ; page++ {
			pagedAccounts, err := client.GetAccounts(ctx, page, pageSize)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
				return err
			}

			if len(pagedAccounts) == 0 {
				break
			}

			if err := ingestAccountsBatch(ctx, connectorID, ingester, pagedAccounts); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedAccounts)), accountsAttrs)

			for _, account := range pagedAccounts {
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

				balancesTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
					Name:      "Fetch balances from client by account",
					Key:       taskNameFetchBalances,
					AccountID: account.ID,
				})

				if err != nil {
					return err
				}
				err = scheduler.Schedule(ctx, balancesTask, models.TaskSchedulerOptions{
					ScheduleOption: models.OPTIONS_RUN_NOW,
					RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
				})
				if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
					return err
				}

				taskRecipients, err := models.EncodeTaskDescriptor(TaskDescriptor{
					Name:      "Fetch recipients from client",
					Key:       taskNameFetchRecipients,
					AccountID: account.ID,
				})
				if err != nil {
					return err
				}

				err = scheduler.Schedule(ctx, taskRecipients, models.TaskSchedulerOptions{
					ScheduleOption: models.OPTIONS_RUN_NOW,
					RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
				})
				if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
					return err
				}
			}

			if len(pagedAccounts) < pageSize {
				break
			}
		}

		return nil
	}
}

func ingestAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
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
				Reference:   account.ID,
				ConnectorID: connectorID,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:   time.Now(),
			Reference:   account.ID,
			ConnectorID: connectorID,
			AccountName: account.Attributes.AccountName,
			Type:        models.AccountTypeInternal,
			RawData:     raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}