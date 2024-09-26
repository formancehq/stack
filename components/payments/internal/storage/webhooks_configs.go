package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/uptrace/bun"
)

type webhookConfig struct {
	bun.BaseModel `bun:"table:webhooks_configs"`

	// Mandatory fields
	Name        string             `bun:"name,pk,type:text,notnull"`
	ConnectorID models.ConnectorID `bun:"connector_id,pk,type:character varying,notnull"`
	URLPath     string             `bun:"url_path,type:text,notnull"`
}

func (s *store) WebhooksConfigsUpsert(ctx context.Context, webhooksConfigs []models.WebhookConfig) error {
	toInsert := fromWebhooksConfigsModels(webhooksConfigs)

	_, err := s.db.NewInsert().
		Model(&toInsert).
		On("CONFLICT (name, connector_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("upsert webhook config", err)
	}

	return nil
}

func (s *store) WebhooksConfigsGet(ctx context.Context, name string, connectorID models.ConnectorID) (*models.WebhookConfig, error) {
	var webhookConfig webhookConfig
	err := s.db.NewSelect().
		Model(&webhookConfig).
		Where("name = ? AND connector_id = ?", name, connectorID).
		Scan(ctx)
	if err != nil {
		return nil, e("get webhook config", err)
	}

	return toWebhookConfigModel(webhookConfig), nil
}

func (s *store) WebhooksConfigsDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*webhookConfig)(nil)).
		Where("connector_id = ?", connectorID).
		Exec(ctx)
	if err != nil {
		return e("delete webhook config", err)
	}

	return nil
}

func fromWebhookConfigModels(from models.WebhookConfig) webhookConfig {
	return webhookConfig{
		Name:        from.Name,
		ConnectorID: from.ConnectorID,
		URLPath:     from.URLPath,
	}
}

func fromWebhooksConfigsModels(from []models.WebhookConfig) []webhookConfig {
	to := make([]webhookConfig, 0, len(from))
	for _, webhookConfig := range from {
		to = append(to, fromWebhookConfigModels(webhookConfig))
	}

	return to
}

func toWebhookConfigModel(from webhookConfig) *models.WebhookConfig {
	return &models.WebhookConfig{
		Name:        from.Name,
		ConnectorID: from.ConnectorID,
		URLPath:     from.URLPath,
	}
}
