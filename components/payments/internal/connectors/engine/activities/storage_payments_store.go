package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StoragePaymentsStore(ctx context.Context, payments []models.Payment) error {
	return a.storage.PaymentsUpsert(ctx, payments)
}

var StoragePaymentsStoreActivity = Activities{}.StoragePaymentsStore

func StoragePaymentsStore(ctx workflow.Context, payments []models.Payment) error {
	return executeActivity(ctx, StoragePaymentsStoreActivity, nil, payments)
}
