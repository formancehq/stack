package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageBankAccountsGet(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	return a.storage.BankAccountsGet(ctx, id, expand)
}

var StorageBankAccountsGetActivity = Activities{}.StorageBankAccountsGet

func StorageBankAccountsGet(ctx workflow.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	var result models.BankAccount
	err := executeActivity(ctx, StorageBankAccountsGetActivity, &result, id, expand)
	return &result, err
}
