package services

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ConnectorsConfigs() plugins.Configs {
	return plugins.GetConfigs()
}

func (s *Service) ConnectorsConfig(ctx context.Context, connectorID models.ConnectorID) (json.RawMessage, error) {
	connector, err := s.storage.ConnectorsGet(ctx, connectorID)
	if err != nil {
		return nil, newStorageError(err, "get connector")
	}

	return connector.Config, nil
}
