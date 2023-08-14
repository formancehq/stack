package stripe

import (
	"context"
	"time"

	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stripe/stripe-go/v72"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	paymentsConnectedAccountsAttrs = append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments_for_connected_account"))
)

func ingestBatch(ctx context.Context, account string, logger logging.Logger, ingester ingestion.Ingester,
	bts []*stripe.BalanceTransaction, commitState TimelineState, tail bool,
) error {
	batch := ingestion.PaymentBatch{}

	for i := range bts {
		batchElement, handled := CreateBatchElement(bts[i], account, !tail)

		if !handled {
			logger.Debugf("Balance transaction type not handled: %s", bts[i].Type)

			continue
		}

		if batchElement.Adjustment == nil && batchElement.Payment == nil {
			continue
		}

		batch = append(batch, batchElement)
	}

	logger.WithFields(map[string]interface{}{
		"state": commitState,
	}).Debugf("updating state")

	err := ingester.IngestPayments(ctx, batch, commitState)
	if err != nil {
		return err
	}

	return nil
}

func ConnectedAccountTask(config Config, account string, client *DefaultClient) func(ctx context.Context, logger logging.Logger,
	ingester ingestion.Ingester, resolver task.StateResolver, metricsRegistry metrics.MetricsRegistry) error {
	return func(ctx context.Context, logger logging.Logger, ingester ingestion.Ingester,
		resolver task.StateResolver, metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Infof("Create new trigger")

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), metric.WithAttributes(paymentsConnectedAccountsAttrs...))
		}()

		trigger := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					if err := ingestBatch(ctx, account, logger, ingester, batch, commitState, tail); err != nil {
						return err
					}
					metricsRegistry.ConnectorObjects().Add(ctx, int64(len(batch)), metric.WithAttributes(paymentsConnectedAccountsAttrs...))

					return nil
				},
				func(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error {
					return nil
				},
			),
			NewTimeline(client.
				ForAccount(account), config.TimelineConfig, task.MustResolveTo(ctx, resolver, TimelineState{})),
			TimelineTriggerTypeTransactions,
		)

		if err := trigger.Fetch(ctx); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, metric.WithAttributes(paymentsConnectedAccountsAttrs...))
			return err
		}

		return nil
	}
}
