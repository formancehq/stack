package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ConnectorsReset(ctx context.Context, connectorID models.ConnectorID) error {
	err := s.engine.ResetConnector(ctx, connectorID)
	return handleEngineErrors(err)
}
