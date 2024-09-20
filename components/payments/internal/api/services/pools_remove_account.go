package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Service) PoolsRemoveAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error {
	err := s.engine.RemoveAccountFromPool(ctx, id, accountID)
	return handleEngineErrors(err)
}
