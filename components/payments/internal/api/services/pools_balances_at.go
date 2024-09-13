package services

import (
	"context"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Service) PoolsBalancesAt(ctx context.Context, poolID uuid.UUID, at time.Time) ([]models.AggregatedBalance, error) {
	pool, err := s.storage.PoolsGet(ctx, poolID)
	if err != nil {
		return nil, newStorageError(err, "getting pool")
	}
	res := make(map[string]*big.Int)
	for _, poolAccount := range pool.PoolAccounts {
		balances, err := s.storage.BalancesGetAt(ctx, poolAccount.AccountID, at)
		if err != nil {
			return nil, newStorageError(err, "getting balances")
		}

		for _, balance := range balances {
			amount, ok := res[balance.Asset]
			if !ok {
				amount = big.NewInt(0)
			}

			amount.Add(amount, balance.Balance)
			res[balance.Asset] = amount
		}
	}

	balances := make([]models.AggregatedBalance, 0, len(res))
	for asset, amount := range res {
		balances = append(balances, models.AggregatedBalance{
			Asset:  asset,
			Amount: amount,
		})
	}

	return balances, nil
}
