package service

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ListBalances(ctx context.Context, query storage.BalanceQuery) ([]*models.Balance, storage.PaginationDetails, error) {
	balances, paginationDetails, err := s.store.ListBalances(ctx, query)
	return balances, paginationDetails, newStorageError(err, "listing balances")
}
