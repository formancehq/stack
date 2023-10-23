package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) IsInstalled(ctx context.Context, provider models.ConnectorProvider) (bool, error) {
	exists, err := s.db.NewSelect().
		Model(&models.Connector{}).
		Where("provider = ?", provider).
		Exists(ctx)
	if err != nil {
		return false, e("find connector", err)
	}

	return exists, nil
}

func (s *Storage) GetConnector(ctx context.Context, provider models.ConnectorProvider) (*models.Connector, error) {
	var connector models.Connector

	err := s.db.NewSelect().
		Model(&connector).
		ColumnExpr("*, pgp_sym_decrypt(config, ?, ?) AS decrypted_config", s.configEncryptionKey, encryptionOptions).
		Where("provider = ?", provider).
		Scan(ctx)
	if err != nil {
		return nil, e("find connector", err)
	}

	return &connector, nil
}
