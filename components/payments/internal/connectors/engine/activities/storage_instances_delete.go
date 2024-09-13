package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageInstancesDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.InstancesDeleteFromConnectorID(ctx, connectorID)
}

var StorageInstancesDeleteActivity = Activities{}.StorageInstancesDelete

func StorageInstancesDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageInstancesDeleteActivity, nil, connectorID)
}
