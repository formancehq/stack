package stripe

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
	"go.opentelemetry.io/otel/attribute"
)

var (
	accountsAttrs = append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))
)

func FetchAccountsTask(config Config, client *DefaultClient) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), accountsAttrs...)
		}()

		tt := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					return nil

				},
				func(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error {
					if err := ingestAccountsBatch(ctx, ingester, batch); err != nil {
						return err
					}
					metricsRegistry.ConnectorObjects().Add(ctx, int64(len(batch)), accountsAttrs...)

					for _, account := range batch {
						transactionsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
							Name:    "Fetch transactions for a specific connected account",
							Key:     taskNameFetchPaymentsForAccounts,
							Account: account.ID,
						})
						if err != nil {
							return errors.Wrap(err, "failed to transform task descriptor")
						}

						err = scheduler.Schedule(ctx, transactionsTask, models.TaskSchedulerOptions{
							ScheduleOption: models.OPTIONS_RUN_NOW,
							Restart:        true,
						})
						if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
							return errors.Wrap(err, "scheduling connected account")
						}

						balancesTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
							Name:    "Fetch balance for a specific connected account",
							Key:     taskNameFetchBalances,
							Account: account.ID,
						})
						if err != nil {
							return errors.Wrap(err, "failed to transform task descriptor")
						}

						err = scheduler.Schedule(ctx, balancesTask, models.TaskSchedulerOptions{
							ScheduleOption: models.OPTIONS_RUN_NOW,
							Restart:        true,
						})
						if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
							return errors.Wrap(err, "scheduling connected account")
						}
					}

					return nil
				},
			),
			NewTimeline(client,
				config.TimelineConfig, task.MustResolveTo(ctx, resolver, TimelineState{})),
			TimelineTriggerTypeAccounts,
		)

		if err := tt.Fetch(ctx); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs...)
			return err
		}

		return nil
	}
}

func ingestAccountsBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	accounts []*stripe.Account,
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
				Provider:  models.ConnectorProviderStripe,
			},
			CreatedAt:    time.Unix(account.Created, 0),
			Reference:    account.ID,
			Provider:     models.ConnectorProviderStripe,
			DefaultAsset: currency.FormatAsset(string(account.DefaultCurrency)),
			Type:         models.AccountTypeInternal,
			RawData:      raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
