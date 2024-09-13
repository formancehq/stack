package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageWebhooksConfigsDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.WebhooksConfigsDeleteFromConnectorID(ctx, connectorID)
}

var StorageWebhooksConfigsDeleteActivity = Activities{}.StorageWebhooksConfigsDelete

func StorageWebhooksConfigsDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageWebhooksConfigsDeleteActivity, nil, connectorID)
}
