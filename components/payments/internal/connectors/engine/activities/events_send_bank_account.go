package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) EventsSendBankAccount(ctx context.Context, bankAccount models.BankAccount) error {
	return a.events.Publish(ctx, a.events.NewEventSavedBankAccounts(bankAccount))
}

var EventsSendBankAccountActivity = Activities{}.EventsSendBankAccount

func EventsSendBankAccount(ctx workflow.Context, bankAccount models.BankAccount) error {
	return executeActivity(ctx, EventsSendBankAccountActivity, nil, bankAccount)
}
