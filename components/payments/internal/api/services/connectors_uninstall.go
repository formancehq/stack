package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ConnectorsUninstall(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.storage.ConnectorsGet(ctx, connectorID)
	if err != nil {
		return newStorageError(err, "get connector")
	}

	return s.engine.UninstallConnector(ctx, connectorID)
}
