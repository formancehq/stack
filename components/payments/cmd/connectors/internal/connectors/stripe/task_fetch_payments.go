package stripe

import (
	"context"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stripe/stripe-go/v72"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func fetchPaymentsTask(config TimelineConfig, client *client.DefaultClient) func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID, resolver task.StateResolver,
	scheduler task.Scheduler, ingester ingestion.Ingester) error {
	return func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID, resolver task.StateResolver,
		scheduler task.Scheduler, ingester ingestion.Ingester,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("stripe.fetchPaymentsTask")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("account", rootAccountReference),
		)

		tt := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					if err := ingestBatch(ctx, connectorID, rootAccountReference, logger, ingester, batch, commitState, tail); err != nil {
						return err
					}

					return nil
				},
				func(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error {
					return nil
				},
				func(ctx context.Context, batch []*stripe.ExternalAccount, commitState TimelineState, tail bool) error {
					return nil
				},
			),
			NewTimeline(client,
				config, task.MustResolveTo(ctx, resolver, TimelineState{})),
			TimelineTriggerTypeTransactions,
		)

		if err := tt.Fetch(ctx); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}
