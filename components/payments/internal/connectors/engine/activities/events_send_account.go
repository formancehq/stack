package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) EventsSendAccount(ctx context.Context, account models.Account) error {
	return a.events.Publish(ctx, a.events.NewEventSavedAccounts(account))
}

var EventsSendAccountActivity = Activities{}.EventsSendAccount

func EventsSendAccount(ctx workflow.Context, account models.Account) error {
	return executeActivity(ctx, EventsSendAccountActivity, nil, account)
}
