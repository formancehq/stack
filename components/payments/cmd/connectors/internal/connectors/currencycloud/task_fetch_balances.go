package currencycloud

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud/client"
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

func taskFetchBalances(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchAccounts)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), balancesAttrs)
		}()

		page := 1
		for {
			if page < 0 {
				break
			}

			pagedBalances, nextPage, err := client.GetBalances(ctx, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
				return err
			}

			page = nextPage

			if err := ingestBalancesBatch(ctx, connectorID, ingester, pagedBalances); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedBalances)), balancesAttrs)
		}

		return nil
	}
}

func ingestBalancesBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	balances []*client.Balance,
) error {
	batch := ingestion.BalanceBatch{}
	for _, balance := range balances {
		var amount big.Float
		_, ok := amount.SetString(balance.Amount)
		if !ok {
			return fmt.Errorf("failed to parse amount %s", balance.Amount)
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(100)).Int(&amountInt)

		now := time.Now()
		batch = append(batch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   balance.AccountID,
				ConnectorID: connectorID,
			},
			Asset:         models.Asset(fmt.Sprintf("%s/2", balance.Currency)),
			Balance:       &amountInt,
			CreatedAt:     now,
			LastUpdatedAt: now,
			ConnectorID:   connectorID,
		})
	}

	if err := ingester.IngestBalances(ctx, batch, true); err != nil {
		return err
	}

	return nil
}
