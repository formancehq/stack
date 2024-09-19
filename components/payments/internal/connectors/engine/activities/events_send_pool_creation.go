package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) EventsSendPoolCreation(ctx context.Context, pool models.Pool) error {
	return a.events.Publish(ctx, a.events.NewEventSavedPool(pool))
}

var EventsSendPoolCreationActivity = Activities{}.EventsSendPoolCreation

func EventsSendPoolCreation(ctx workflow.Context, pool models.Pool) error {
	return executeActivity(ctx, EventsSendPoolCreationActivity, nil, pool)
}
