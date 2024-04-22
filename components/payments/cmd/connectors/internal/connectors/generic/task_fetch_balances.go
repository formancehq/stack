package generic

import (
	"context"
	"fmt"
	"math/big"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/generic/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/genericclient"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

func taskFetchBalances(client *client.Client, config *Config, accountID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"generic.taskFetchBalances",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		balances, err := client.GetBalances(ctx, accountID)
		if err != nil {
			// retryable error already handled by the client
			otel.RecordError(span, err)
			return err
		}

		if err := ingestBalancesBatch(ctx, connectorID, ingester, accountID, balances); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func ingestBalancesBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	accountID string,
	balances *genericclient.Balances,
) error {
	if balances == nil {
		return nil
	}

	balancesBatch := make([]*models.Balance, 0, len(balances.Balances))
	for _, balance := range balances.Balances {
		var amount big.Int
		_, ok := amount.SetString(balance.Amount, 10)
		if !ok {
			return fmt.Errorf("failed to parse amount: %s", balance.Amount)
		}

		balancesBatch = append(balancesBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   accountID,
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Currency),
			Balance:       &amount,
			CreatedAt:     balances.At,
			LastUpdatedAt: balances.At,
			ConnectorID:   connectorID,
		})
	}

	if err := ingester.IngestBalances(ctx, ingestion.BalanceBatch(balancesBatch), false); err != nil {
		return err
	}

	return nil
}
