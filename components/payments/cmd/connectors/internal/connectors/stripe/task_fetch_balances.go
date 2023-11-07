package stripe

import (
	"context"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	balancesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "balances"))...)
)

func balancesTask(account string, client *client.DefaultClient) func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID,
	ingester ingestion.Ingester, resolver task.StateResolver, metricsRegistry metrics.MetricsRegistry) error {
	return func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID, ingester ingestion.Ingester,
		resolver task.StateResolver, metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Infof("Create new balances trigger for account %s", account)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), balancesAttrs)
		}()

		balances, err := client.ForAccount(account).Balance(ctx)
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
			return err
		}

		batch := ingestion.BalanceBatch{}
		for _, balance := range balances.Available {
			timestamp := time.Now()
			batch = append(batch, &models.Balance{
				AccountID: models.AccountID{
					Reference:   account,
					ConnectorID: connectorID,
				},
				Asset:         currency.FormatAsset(string(balance.Currency)),
				Balance:       big.NewInt(balance.Value),
				CreatedAt:     timestamp,
				LastUpdatedAt: timestamp,
				ConnectorID:   connectorID,
			})
		}

		if err := ingester.IngestBalances(ctx, batch, false); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
			return err
		}
		metricsRegistry.ConnectorObjects().Add(ctx, int64(len(balances.Available)), balancesAttrs)

		return nil
	}
}
