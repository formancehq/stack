package service

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func (s *Service) ListBalances(ctx context.Context, q storage.ListBalancesQuery) (*api.Cursor[models.Balance], error) {
	cursor, err := s.store.ListBalances(ctx, q)
	return cursor, newStorageError(err, "listing balances")
}
