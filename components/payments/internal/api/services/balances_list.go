package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) BalancesList(ctx context.Context, query storage.ListBalancesQuery) (*bunpaginate.Cursor[models.Balance], error) {
	return s.storage.BalancesList(ctx, query)
}