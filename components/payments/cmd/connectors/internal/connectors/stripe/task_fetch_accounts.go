package stripe

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
	"go.opentelemetry.io/otel/attribute"
)

const (
	rootAccountReference = "root"
)

func fetchAccountsTask(config TimelineConfig, client *client.DefaultClient) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"stripe.fetchAccountsTask",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		// Register root account.
		if err := registerRootAccount(ctx, connectorID, ingester, scheduler); err != nil {
			otel.RecordError(span, err)
			return err
		}

		tt := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					return nil

				},
				func(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error {
					if err := ingestAccountsBatch(ctx, connectorID, ingester, batch); err != nil {
						return err
					}

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
							RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
						})
						if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
							return errors.Wrap(err, "scheduling connected account")
						}

						balanceTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
							Name:    "Fetch balance for a specific connected account",
							Key:     taskNameFetchBalances,
							Account: account.ID,
						})
						if err != nil {
							return errors.Wrap(err, "failed to transform task descriptor")
						}

						err = scheduler.Schedule(ctx, balanceTask, models.TaskSchedulerOptions{
							ScheduleOption: models.OPTIONS_RUN_NOW,
							RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
						})
						if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
							return errors.Wrap(err, "scheduling connected account")
						}

						externalAccountsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
							Name:    "Fetch external account for a specific connected account",
							Key:     taskNameFetchExternalAccounts,
							Account: account.ID,
						})
						if err != nil {
							return errors.Wrap(err, "failed to transform task descriptor")
						}

						err = scheduler.Schedule(ctx, externalAccountsTask, models.TaskSchedulerOptions{
							ScheduleOption: models.OPTIONS_RUN_NOW,
							RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
						})
						if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
							return errors.Wrap(err, "scheduling connected account")
						}
					}

					return nil
				},
				func(ctx context.Context, batch []*stripe.ExternalAccount, commitState TimelineState, tail bool) error {
					return nil

				},
			),
			NewTimeline(client,
				config, task.MustResolveTo(ctx, resolver, TimelineState{})),
			TimelineTriggerTypeAccounts,
		)

		if err := tt.Fetch(ctx); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func registerRootAccount(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
) error {
	if err := ingester.IngestAccounts(ctx, ingestion.AccountBatch{
		{
			ID: models.AccountID{
				Reference:   rootAccountReference,
				ConnectorID: connectorID,
			},
			ConnectorID: connectorID,
			CreatedAt:   time.Now().UTC(),
			Reference:   rootAccountReference,
			Type:        models.AccountTypeInternal,
		},
	}); err != nil {
		return err
	}
	balanceTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name:    "Fetch balance for the root account",
		Key:     taskNameFetchBalances,
		Account: rootAccountReference,
	})
	if err != nil {
		return errors.Wrap(err, "failed to transform task descriptor")
	}
	err = scheduler.Schedule(ctx, balanceTask, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
	})
	if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
		return errors.Wrap(err, "scheduling connected account")
	}

	return nil
}

func ingestAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	accounts []*stripe.Account,
) error {
	batch := ingestion.AccountBatch{}
	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		metadata := make(map[string]string)
		for k, v := range account.Metadata {
			metadata[k] = v
		}

		batch = append(batch, &models.Account{
			ID: models.AccountID{
				Reference:   account.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:    time.Unix(account.Created, 0).UTC(),
			Reference:    account.ID,
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, string(account.DefaultCurrency)),
			Type:         models.AccountTypeInternal,
			RawData:      raw,
			Metadata:     metadata,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
