package connectors_manager

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
)

type Store interface {
	ListConnectors(ctx context.Context) ([]*models.Connector, error)
	IsInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error)
	IsInstalledByConnectorName(ctx context.Context, name string) (bool, error)
	Install(ctx context.Context, connector *models.Connector, config json.RawMessage) error
	Uninstall(ctx context.Context, connectorID models.ConnectorID) error
	UpdateConfig(ctx context.Context, connectorID models.ConnectorID, config json.RawMessage) error
	GetConnector(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error)
	CreateWebhook(ctx context.Context, webhook *models.Webhook) error
}
