package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Service) BankAccountsGet(ctx context.Context, id uuid.UUID) (*models.BankAccount, error) {
	return s.storage.BankAccountsGet(ctx, id, true)
}
