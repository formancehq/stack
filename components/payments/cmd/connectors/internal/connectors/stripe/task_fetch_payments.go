package stripe

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stripe/stripe-go/v72"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	paymentsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments_for_connected_account"))...)
)

func fetchPaymentsTask(config TimelineConfig, client *client.DefaultClient) func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID, resolver task.StateResolver,
	scheduler task.Scheduler, ingester ingestion.Ingester, metricsRegistry metrics.MetricsRegistry) error {
	return func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID, resolver task.StateResolver,
		scheduler task.Scheduler, ingester ingestion.Ingester, metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsAttrs)
		}()

		tt := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					if err := ingestBatch(ctx, connectorID, "", logger, ingester, batch, commitState, tail); err != nil {
						return err
					}
					metricsRegistry.ConnectorObjects().Add(ctx, int64(len(batch)), paymentsAttrs)

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
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
			return err
		}

		return nil
	}
}
