package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) PoolsCreate(ctx context.Context, pool models.Pool) error {
	err := s.engine.CreatePool(ctx, pool)
	return handleEngineErrors(err)
}
