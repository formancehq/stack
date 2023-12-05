package moneycorp

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
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

func taskFetchBalances(logger logging.Logger, client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info("Fetching transactions for account", accountID)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), balancesAttrs)
		}()

		balances, err := client.GetAccountBalances(ctx, accountID)
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
			return err
		}

		if err := ingestBalancesBatch(ctx, connectorID, ingester, accountID, balances); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
			return err
		}
		metricsRegistry.ConnectorObjects().Add(ctx, int64(len(balances)), balancesAttrs)

		return nil
	}
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
