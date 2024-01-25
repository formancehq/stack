package moneycorp

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

func taskFetchBalances(client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"moneycorp.taskFetchBalances",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("accountID", accountID),
		)
		defer span.End()

		if err := fetchBalances(ctx, client, accountID, connectorID, ingester); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchBalances(
	ctx context.Context,
	client *client.Client,
	accountID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
) error {
	balances, err := client.GetAccountBalances(ctx, accountID)
	if err != nil {
		return err
	}

	if err := ingestBalancesBatch(ctx, connectorID, ingester, accountID, balances); err != nil {
		return err
	}

	return nil
}

func ingestBalancesBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	accountID string,
	balances []*client.Balance,
) error {
	batch := ingestion.BalanceBatch{}
	for _, balance := range balances {
		var amount big.Float
		_, ok := amount.SetString(balance.Attributes.AvailableBalance.String())
		if !ok {
			return fmt.Errorf("failed to parse amount %s", balance.Attributes.AvailableBalance.String())
		}

		precision, err := currency.GetPrecision(supportedCurrenciesWithDecimal, balance.Attributes.CurrencyCode)
		if err != nil {
			return err
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

		now := time.Now()
		batch = append(batch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   accountID,
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Attributes.CurrencyCode),
			Balance:       &amountInt,
			CreatedAt:     now,
			LastUpdatedAt: now,
			ConnectorID:   connectorID,
		})
	}

	return ingester.IngestBalances(ctx, batch, false)
}
