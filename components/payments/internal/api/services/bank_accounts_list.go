package services

import (
	"context"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
)

func (s *Service) BankAccountsList(ctx context.Context, query storage.ListBankAccountsQuery) (*bunpaginate.Cursor[models.BankAccount], error) {
	return s.storage.BankAccountsList(ctx, query)
}
