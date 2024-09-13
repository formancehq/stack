package services

import (
	"context"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
)

func (s *Service) AccountsList(ctx context.Context, query storage.ListAccountsQuery) (*bunpaginate.Cursor[models.Account], error) {
	return s.storage.AccountsList(ctx, query)
}
