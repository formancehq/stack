package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageInstancesStore(ctx context.Context, instance models.Instance) error {
	return a.storage.InstancesUpsert(ctx, instance)
}

var StorageInstancesStoreActivity = Activities{}.StorageInstancesStore

func StorageInstancesStore(ctx workflow.Context, instance models.Instance) error {
	return executeActivity(ctx, StorageInstancesStoreActivity, nil, instance)
}
