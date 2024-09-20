package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) EventsSendConnectorReset(ctx context.Context, connectorID models.ConnectorID) error {
	return a.events.Publish(ctx, a.events.NewEventResetConnector(connectorID))
}

var EventsSendConnectorResetActivity = Activities{}.EventsSendConnectorReset

func EventsSendConnectorReset(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, EventsSendConnectorResetActivity, nil, connectorID)
}
