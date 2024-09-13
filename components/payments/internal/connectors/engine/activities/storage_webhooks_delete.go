package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageWebhooksDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.WebhooksDeleteFromConnectorID(ctx, connectorID)
}

var StorageWebhooksDeleteActivity = Activities{}.StorageWebhooksDelete

func StorageWebhooksDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageWebhooksDeleteActivity, nil, connectorID)
}
