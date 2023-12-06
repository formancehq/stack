package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) ListConnectors(ctx context.Context) ([]*models.Connector, error) {
	var connectors []*models.Connector

	err := s.db.NewSelect().
		Model(&connectors).
		ColumnExpr("*, pgp_sym_decrypt(config, ?, ?) AS decrypted_config", s.configEncryptionKey, encryptionOptions).
		Scan(ctx)
	if err != nil {
		return nil, e("list connectors", err)
	}

	return connectors, nil
}

func (s *Storage) ListConnectorsByProvider(ctx context.Context, provider models.ConnectorProvider) ([]*models.Connector, error) {
	var connectors []*models.Connector

	err := s.db.NewSelect().
		Model(&connectors).
		ColumnExpr("*, pgp_sym_decrypt(config, ?, ?) AS decrypted_config", s.configEncryptionKey, encryptionOptions).
		Where("provider = ?", provider).
		Scan(ctx)
	if err != nil {
		return nil, e("list connectors", err)
	}

	return connectors, nil
}

func (s *Storage) GetConfig(ctx context.Context, connectorID models.ConnectorID, destination any) error {
	var connector models.Connector

	err := s.db.NewSelect().
		Model(&connector).
		ColumnExpr("pgp_sym_decrypt(config, ?, ?) AS decrypted_config", s.configEncryptionKey, encryptionOptions).
		Where("id = ?", connectorID).
		Scan(ctx)
	if err != nil {
		return e(fmt.Sprintf("failed to get config for connector %s", connectorID), err)
	}

	err = json.Unmarshal(connector.Config, destination)
	if err != nil {
		return e(fmt.Sprintf("failed to unmarshal config for connector %s", connectorID), err)
	}

	return nil
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

func (s *Storage) Install(ctx context.Context, connector *models.Connector, config json.RawMessage) error {
	_, err := s.db.NewInsert().Model(connector).Exec(ctx)
	if err != nil {
		return e("install connector", err)
	}

	return s.UpdateConfig(ctx, connector.ID, config)
}

func (s *Storage) Uninstall(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model(&models.Connector{}).
		Where("id = ?", connectorID).
		Exec(ctx)
	if err != nil {
		return e("uninstall connector", err)
	}

	return nil
}

func (s *Storage) UpdateConfig(ctx context.Context, connectorID models.ConnectorID, config json.RawMessage) error {
	_, err := s.db.NewUpdate().
		Model(&models.Connector{}).
		Set("config = pgp_sym_encrypt(?::TEXT, ?, ?)", config, s.configEncryptionKey, encryptionOptions).
		Where("id = ?", connectorID). // Connector name is unique
		Exec(ctx)
	if err != nil {
		return e("update connector config", err)
	}

	return nil
}

func (s *Storage) GetConnector(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error) {
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
