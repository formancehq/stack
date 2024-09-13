package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageInstancesUpdate(ctx context.Context, instance models.Instance) error {
	return a.storage.InstancesUpdate(ctx, instance)
}

var StorageInstancesUpdateActivity = Activities{}.StorageInstancesUpdate

func StorageInstancesUpdate(ctx workflow.Context, instance models.Instance) error {
	return executeActivity(ctx, StorageInstancesUpdateActivity, nil, instance)
}
