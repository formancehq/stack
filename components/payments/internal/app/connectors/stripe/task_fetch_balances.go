package stripe

import (
	"context"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func BalancesTask(config Config, account string, client *DefaultClient) func(ctx context.Context, logger logging.Logger,
	ingester ingestion.Ingester, resolver task.StateResolver) error {
	return func(ctx context.Context, logger logging.Logger, ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		logger.Infof("Create new balances trigger for account %s", account)
		balances, err := client.ForAccount(account).Balance(ctx)
		if err != nil {
			return err
		}

		batch := ingestion.BalanceBatch{}
		now := time.Now()
		for _, balance := range balances.Available {
			batch = append(batch, &models.Balance{
				AccountID: models.AccountID{
					Reference: account,
					Provider:  models.ConnectorProviderStripe,
				},
				Currency:      currency.FormatAsset(string(balance.Currency)).String(),
				Balance:       big.NewInt(balance.Value),
				CreatedAt:     now,
				LastUpdatedAt: now,
			})
		}

		if err := ingester.IngestBalances(ctx, batch, false); err != nil {
			return err
		}

		return nil
	}
}
