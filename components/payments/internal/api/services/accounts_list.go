package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) AccountsList(ctx context.Context, query storage.ListAccountsQuery) (*bunpaginate.Cursor[models.Account], error) {
	return s.storage.AccountsList(ctx, query)
}
