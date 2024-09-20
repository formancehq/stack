package stripe

import (
	"context"
	"math/big"
	"time"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

func balanceTask(account string, client *client.DefaultClient) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"stripe.balanceTask",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("account", account),
		)
		defer span.End()

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
