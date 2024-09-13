package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageBankAccountsAddRelatedAccount(ctx context.Context, relatedAccount models.BankAccountRelatedAccount) error {
	return a.storage.BankAccountsAddRelatedAccount(ctx, relatedAccount)
}

var StorageBankAccountsAddRelatedAccountActivity = Activities{}.StorageBankAccountsAddRelatedAccount

func StorageBankAccountsAddRelatedAccount(ctx workflow.Context, relatedAccount models.BankAccountRelatedAccount) error {
	return executeActivity(ctx, StorageBankAccountsAddRelatedAccountActivity, nil, relatedAccount)
}
