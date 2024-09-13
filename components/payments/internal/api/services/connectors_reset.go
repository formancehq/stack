package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ConnectorsReset(ctx context.Context, connectorID models.ConnectorID) error {
	config, err := s.ConnectorsConfig(ctx, connectorID)
	if err != nil {
		return err
	}

	if err := s.engine.UninstallConnector(ctx, connectorID); err != nil {
		return err
	}

	_, err = s.engine.InstallConnector(ctx, connectorID.Provider, config)
	return err
}
