package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageBalancesDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.BalancesDeleteForConnectorID(ctx, connectorID)
}

var StorageBalancesDeleteActivity = Activities{}.StorageBalancesDelete

func StorageBalancesDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageBalancesDeleteActivity, nil, connectorID)
}
