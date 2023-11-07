package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) IsInstalledByConnectorName(ctx context.Context, name string) (bool, error) {
	exists, err := s.db.NewSelect().
		Model(&models.Connector{}).
		Where("name = ?", name).
		Exists(ctx)
	if err != nil {
		return false, e("find connector", err)
	}

	return exists, nil
}

func (s *Storage) IsInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error) {
	exists, err := s.db.NewSelect().
		Model(&models.Connector{}).
		Where("id = ?", connectorID).
		Exists(ctx)
	if err != nil {
		return false, e("find connector", err)
	}

	return exists, nil
}

func (s *Storage) GetConnectorByConnectorID(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error) {
	var connector models.Connector

	err := s.db.NewSelect().
		Model(&connector).
		ColumnExpr("*, pgp_sym_decrypt(config, ?, ?) AS decrypted_config", s.configEncryptionKey, encryptionOptions).
		Where("id = ?", connectorID).
		Scan(ctx)
	if err != nil {
		return nil, e("find connector", err)
	}

	return &connector, nil
}
