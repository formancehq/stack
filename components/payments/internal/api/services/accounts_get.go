package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) AccountsGet(ctx context.Context, id models.AccountID) (*models.Account, error) {
	return s.storage.AccountsGet(ctx, id)
}
