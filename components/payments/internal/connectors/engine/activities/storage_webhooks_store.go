package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageWebhooksStore(ctx context.Context, webhook models.Webhook) error {
	return a.storage.WebhooksInsert(ctx, webhook)
}

var StorageWebhooksStoreActivity = Activities{}.StorageWebhooksStore

func StorageWebhooksStore(ctx workflow.Context, webhook models.Webhook) error {
	return executeActivity(ctx, StorageWebhooksStoreActivity, nil, webhook)
}
