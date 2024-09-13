package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageTasksTreeDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.TasksDeleteFromConnectorID(ctx, connectorID)
}

var StorageTasksTreeDeleteActivity = Activities{}.StorageTasksTreeDelete

func StorageTasksTreeDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageTasksTreeDeleteActivity, nil, connectorID)
}
