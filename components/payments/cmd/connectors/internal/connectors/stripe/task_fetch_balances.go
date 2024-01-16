package stripe

import (
	"context"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func balanceTask(account string, client *client.DefaultClient) func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID,
	ingester ingestion.Ingester, resolver task.StateResolver) error {
	return func(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID, ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("stripe.balanceTask")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("account", account),
		)

		stripeAccount := account
		if account == rootAccountReference {
			// special case for root account
			stripeAccount = ""
		}

		balances, err := client.ForAccount(stripeAccount).Balance(ctx)
		if err != nil {
			otel.RecordError(span, err)
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
				Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, string(balance.Currency)),
				Balance:       big.NewInt(balance.Value),
				CreatedAt:     timestamp,
				LastUpdatedAt: timestamp,
				ConnectorID:   connectorID,
			})
		}

		if err := ingester.IngestBalances(ctx, batch, false); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}
