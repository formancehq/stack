package services

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ConnectorsInstall(ctx context.Context, provider string, config json.RawMessage) (models.ConnectorID, error) {
	connectorID, err := s.engine.InstallConnector(ctx, provider, config)
	return connectorID, handleEngineErrors(err)
}
