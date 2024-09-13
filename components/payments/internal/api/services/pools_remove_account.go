package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Service) PoolsRemoveAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error {
	return s.storage.PoolsRemoveAccount(ctx, id, accountID)
}
