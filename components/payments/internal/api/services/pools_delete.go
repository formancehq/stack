package services

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) PoolsDelete(ctx context.Context, id uuid.UUID) error {
	return s.storage.PoolsDelete(ctx, id)
}
