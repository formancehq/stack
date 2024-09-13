package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Service) PoolsGet(ctx context.Context, id uuid.UUID) (*models.Pool, error) {
	return s.storage.PoolsGet(ctx, id)
}
