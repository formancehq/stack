package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) EventsSendBalance(ctx context.Context, balance models.Balance) error {
	return a.events.Publish(ctx, a.events.NewEventSavedBalances(balance))
}

var EventsSendBalanceActivity = Activities{}.EventsSendBalance

func EventsSendBalance(ctx workflow.Context, balance models.Balance) error {
	return executeActivity(ctx, EventsSendBalanceActivity, nil, balance)
}
