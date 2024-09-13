package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageSchedulesDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.SchedulesDeleteFromConnectorID(ctx, connectorID)
}

var StorageSchedulesDeleteActivity = Activities{}.StorageSchedulesDelete

func StorageSchedulesDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageSchedulesDeleteActivity, nil, connectorID)
}
