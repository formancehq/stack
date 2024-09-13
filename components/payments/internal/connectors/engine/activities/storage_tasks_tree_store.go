package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

type TasksTreeStoreRequest struct {
	ConnectorID models.ConnectorID
	Workflow    models.Tasks
}

func (a Activities) StorageTasksTreeStore(ctx context.Context, request TasksTreeStoreRequest) error {
	return a.storage.TasksUpsert(ctx, request.ConnectorID, request.Workflow)
}

var StorageTasksTreeStoreActivity = Activities{}.StorageTasksTreeStore

func StorageTasksTreeStore(ctx workflow.Context, connectorID models.ConnectorID, workflow models.Tasks) error {
	return executeActivity(ctx, StorageTasksTreeStoreActivity, nil, TasksTreeStoreRequest{
		ConnectorID: connectorID,
		Workflow:    workflow,
	})
}
