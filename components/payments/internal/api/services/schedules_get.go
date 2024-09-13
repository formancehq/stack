package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) SchedulesGet(ctx context.Context, id string, connectorID models.ConnectorID) (*models.Schedule, error) {
	return s.storage.SchedulesGet(ctx, id, connectorID)
}
