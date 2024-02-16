package service

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pkg/errors"
)

func (s *Service) ListAccounts(ctx context.Context, q storage.ListAccountsQuery) (*api.Cursor[models.Account], error) {
	cursor, err := s.store.ListAccounts(ctx, q)
	return cursor, newStorageError(err, "listing accounts")
}

func (s *Service) GetAccount(
	ctx context.Context,
	accountID string,
) (*models.Account, error) {
	_, err := models.AccountIDFromString(accountID)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	account, err := s.store.GetAccount(ctx, accountID)
	return account, newStorageError(err, "getting account")
}
