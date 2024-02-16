package service

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

func (s *Service) ListBankAccounts(ctx context.Context, q storage.ListBankAccountQuery) (*api.Cursor[models.BankAccount], error) {
	cursor, err := s.store.ListBankAccounts(ctx, q)
	return cursor, newStorageError(err, "listing bank accounts")
}

func (s *Service) GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	account, err := s.store.GetBankAccount(ctx, id, expand)
	return account, newStorageError(err, "getting bank account")
}
