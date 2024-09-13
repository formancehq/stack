package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageAccountsDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.AccountsDeleteFromConnectorID(ctx, connectorID)
}

var StorageAccountsDeleteActivity = Activities{}.StorageAccountsDelete

func StorageAccountsDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageAccountsDeleteActivity, nil, connectorID)
}
