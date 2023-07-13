package stripe

import (
	"context"

	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
)

func ingest(
	ctx context.Context,
	logger logging.Logger,
	scheduler task.Scheduler,
	ingester ingestion.Ingester,
	bts []*stripe.BalanceTransaction,
	commitState TimelineState,
	tail bool,
) error {
	err := ingestBatch(ctx, "", logger, ingester, bts, commitState, tail)
	if err != nil {
		return err
	}

	connectedAccounts := make([]string, 0)

	for _, bt := range bts {
		if bt.Type == stripe.BalanceTransactionTypeTransfer {
			connectedAccounts = append(connectedAccounts, bt.Source.Transfer.Destination.ID)
		}
	}

	for _, connectedAccount := range connectedAccounts {
		descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:    "Fetch balance transactions for a specific connected account",
			Key:     taskNameFetchPaymentsForAccounts,
			Account: connectedAccount,
		})
		if err != nil {
			return errors.Wrap(err, "failed to transform task descriptor")
		}

		err = scheduler.Schedule(ctx, descriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			Restart:        true,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return errors.Wrap(err, "scheduling connected account")
		}
	}

	return nil
}

func FetchPaymentsTask(config Config, client *DefaultClient) func(ctx context.Context, logger logging.Logger, resolver task.StateResolver,
	scheduler task.Scheduler, ingester ingestion.Ingester) error {
	return func(ctx context.Context, logger logging.Logger, resolver task.StateResolver,
		scheduler task.Scheduler, ingester ingestion.Ingester,
	) error {
		tt := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					return ingest(ctx, logger, scheduler, ingester, batch, commitState, tail)

				},
				func(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error {
					return nil
				},
			),
			NewTimeline(client,
				config.TimelineConfig, task.MustResolveTo(ctx, resolver, TimelineState{})),
			TimelineTriggerTypeTransactions,
		)

		return tt.Fetch(ctx)
	}
}
