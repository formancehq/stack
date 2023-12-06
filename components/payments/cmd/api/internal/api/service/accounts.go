package service

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) ListAccounts(
	ctx context.Context,
	pagination storage.PaginatorQuery,
) ([]*models.Account, storage.PaginationDetails, error) {
	accounts, paginationDetails, err := s.store.ListAccounts(ctx, pagination)
	return accounts, paginationDetails, newStorageError(err, "listing accounts")
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
