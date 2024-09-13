package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) BankAccountsCreate(ctx context.Context, bankAccount models.BankAccount) error {
	return s.storage.BankAccountsUpsert(ctx, bankAccount)
}
