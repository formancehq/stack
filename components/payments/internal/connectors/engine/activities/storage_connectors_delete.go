package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageConnectorsDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.ConnectorsUninstall(ctx, connectorID)
}

var StorageConnectorsDeleteActivity = Activities{}.StorageConnectorsDelete

func StorageConnectorsDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageConnectorsDeleteActivity, nil, connectorID)
}
