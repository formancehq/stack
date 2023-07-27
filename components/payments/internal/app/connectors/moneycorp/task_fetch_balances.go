package moneycorp

import (
	"context"
	"math"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/connectors/moneycorp/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchBalances(logger logging.Logger, client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
	) error {
		logger.Info("Fetching transactions for account", accountID)

		balances, err := client.GetAccountBalances(ctx, accountID)
		if err != nil {
			return err
		}

		if err := ingestBalancesBatch(ctx, ingester, accountID, balances); err != nil {
			return err
		}

		return nil
	}
}

func ingestBalancesBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	accountID string,
	balances []*client.Balance,
) error {
	batch := ingestion.BalanceBatch{}
	now := time.Now()
	for _, balance := range balances {
		batch = append(batch, &models.Balance{
			AccountID: models.AccountID{
				Reference: accountID,
				Provider:  models.ConnectorProviderMoneycorp,
			},
			Currency:      currency.FormatAsset(balance.Attributes.CurrencyCode).String(),
			Balance:       int64(balance.Attributes.AvailableBalance * math.Pow(10, float64(currency.GetPrecision(balance.Attributes.CurrencyCode)))),
			CreatedAt:     now,
			LastUpdatedAt: now,
		})
	}

	return ingester.IngestBalances(ctx, batch, false)
}
