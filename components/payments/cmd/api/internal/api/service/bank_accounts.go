package service

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Service) ListBankAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.BankAccount, storage.PaginationDetails, error) {
	accounts, paginationDetails, err := s.store.ListBankAccounts(ctx, pagination)
	return accounts, paginationDetails, newStorageError(err, "listing bank accounts")
}

func (s *Service) GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	account, err := s.store.GetBankAccount(ctx, id, expand)
	return account, newStorageError(err, "getting bank account")
}
