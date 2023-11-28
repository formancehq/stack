package service

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ListConnectors(ctx context.Context) ([]*models.Connector, error) {
	connectors, err := s.store.ListConnectors(ctx)
	return connectors, newStorageError(err, "listing connectors")
}
