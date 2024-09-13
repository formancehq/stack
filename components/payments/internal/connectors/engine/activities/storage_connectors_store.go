package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageConnectorsStore(ctx context.Context, connector models.Connector) error {
	return a.storage.ConnectorsInstall(ctx, connector)
}

var StorageConnectorsStoreActivity = Activities{}.StorageConnectorsStore

func StorageConnectorsStore(ctx workflow.Context, connector models.Connector) error {
	return executeActivity(ctx, StorageConnectorsStoreActivity, nil, connector)
}
