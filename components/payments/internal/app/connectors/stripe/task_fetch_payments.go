package stripe

import (
	"context"

	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stripe/stripe-go/v72"
)

func FetchPaymentsTask(config Config, client *DefaultClient) func(ctx context.Context, logger logging.Logger, resolver task.StateResolver,
	scheduler task.Scheduler, ingester ingestion.Ingester) error {
	return func(ctx context.Context, logger logging.Logger, resolver task.StateResolver,
		scheduler task.Scheduler, ingester ingestion.Ingester,
	) error {
		tt := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					return ingestBatch(ctx, "", logger, ingester, batch, commitState, tail)
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
