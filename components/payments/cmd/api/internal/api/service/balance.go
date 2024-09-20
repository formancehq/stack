package service

import (
	"context"

	"github.com/formancehq/go-libs/bun/bunpaginate"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ListBalances(ctx context.Context, q storage.ListBalancesQuery) (*bunpaginate.Cursor[models.Balance], error) {
	cursor, err := s.store.ListBalances(ctx, q)
	return cursor, newStorageError(err, "listing balances")
}
