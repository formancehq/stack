package services

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) PoolsDelete(ctx context.Context, id uuid.UUID) error {
	err := s.engine.DeletePool(ctx, id)
	return handleEngineErrors(err)
}
