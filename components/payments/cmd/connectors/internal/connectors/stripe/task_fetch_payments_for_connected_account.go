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
	paymentsConnectedAccountsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments_for_connected_account"))...)
)

func ingestBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	account string,
	logger logging.Logger,
	ingester ingestion.Ingester,
	bts []*stripe.BalanceTransaction,
	commitState TimelineState,
	tail bool,
) error {
	batch := ingestion.PaymentBatch{}

	for i := range bts {
		batchElement, handled := createBatchElement(connectorID, bts[i], account, !tail)

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

	err := ingester.IngestPayments(ctx, connectorID, batch, commitState)
	if err != nil {
		return err
	}

	return nil
}

func connectedAccountTask(config TimelineConfig, account string, client *client.DefaultClient) func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID,
	ingester ingestion.Ingester, resolver task.StateResolver, metricsRegistry metrics.MetricsRegistry) error {
	return func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID, ingester ingestion.Ingester,
		resolver task.StateResolver, metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Infof("Create new trigger")

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsConnectedAccountsAttrs)
		}()

		trigger := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					if err := ingestBatch(ctx, connectorID, account, logger, ingester, batch, commitState, tail); err != nil {
						return err
					}
					metricsRegistry.ConnectorObjects().Add(ctx, int64(len(batch)), paymentsConnectedAccountsAttrs)

					return nil
				},
				func(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error {
					return nil
				},
				func(ctx context.Context, batch []*stripe.ExternalAccount, commitState TimelineState, tail bool) error {
					return nil
				},
			),
			NewTimeline(client.
				ForAccount(account), config, task.MustResolveTo(ctx, resolver, TimelineState{})),
			TimelineTriggerTypeTransactions,
		)

		if err := trigger.Fetch(ctx); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsConnectedAccountsAttrs)
			return err
		}

		return nil
	}
}
