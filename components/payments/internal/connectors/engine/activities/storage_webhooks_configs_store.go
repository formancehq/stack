package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageWebhooksConfigsStore(ctx context.Context, configs []models.WebhookConfig) error {
	return a.storage.WebhooksConfigsUpsert(ctx, configs)
}

var StorageWebhooksConfigsStoreActivity = Activities{}.StorageWebhooksConfigsStore

func StorageWebhooksConfigsStore(ctx workflow.Context, configs []models.WebhookConfig) error {
	return executeActivity(ctx, StorageWebhooksConfigsStoreActivity, nil, configs)
}
