package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStatesDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.StatesDeleteFromConnectorID(ctx, connectorID)
}

var StorageStatesDeleteActivity = Activities{}.StorageStatesDelete

func StorageStatesDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageStatesDeleteActivity, nil, connectorID)
}
