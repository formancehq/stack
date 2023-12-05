package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) IsConnectorInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error) {
	exists, err := s.db.NewSelect().
		Model(&models.Connector{}).
		Where("id = ?", connectorID).
		Exists(ctx)
	if err != nil {
		return false, e("find connector", err)
	}

	return exists, nil
}
