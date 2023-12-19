package currencycloud

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchBalances(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		logger.Info(taskNameFetchAccounts)

		page := 1
		for {
			if page < 0 {
				break
			}

			pagedBalances, nextPage, err := client.GetBalances(ctx, page)
			if err != nil {
				return err
			}

			page = nextPage

			if err := ingestBalancesBatch(ctx, connectorID, ingester, pagedBalances); err != nil {
				return err
			}
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
		// No need to check if the currency is supported for accounts and balances.
		precision := supportedCurrenciesWithDecimal[balance.Currency]

		var amount big.Float
		_, ok := amount.SetString(balance.Amount)
		if !ok {
			return fmt.Errorf("failed to parse amount %s", balance.Amount)
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

		now := time.Now()
		batch = append(batch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   balance.AccountID,
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Currency),
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
