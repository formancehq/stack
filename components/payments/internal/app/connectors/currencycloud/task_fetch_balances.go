package currencycloud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currencycloud/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchBalances(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
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

			if err := ingestBalancesBatch(ctx, ingester, pagedBalances); err != nil {
				return err
			}
		}

		return nil
	}
}

func ingestBalancesBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	balances []*client.Balance,
) error {
	batch := ingestion.BalanceBatch{}
	now := time.Now()
	for _, balance := range balances {
		amount, err := strconv.ParseFloat(balance.Amount, 64)
		if err != nil {
			return err
		}

		batch = append(batch, &models.Balance{
			AccountID: models.AccountID{
				Reference: balance.AccountID,
				Provider:  models.ConnectorProviderCurrencyCloud,
			},
			Currency:      fmt.Sprintf("%s/2", balance.Currency),
			Balance:       int64(amount * 100),
			CreatedAt:     now,
			LastUpdatedAt: now,
		})
	}

	if err := ingester.IngestBalances(ctx, batch, true); err != nil {
		return err
	}

	return nil
}
