package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageAccountsStore(ctx context.Context, accounts []models.Account) error {
	return a.storage.AccountsUpsert(ctx, accounts)
}

var StorageAccountsStoreActivity = Activities{}.StorageAccountsStore

func StorageAccountsStore(ctx workflow.Context, accounts []models.Account) error {
	return executeActivity(ctx, StorageAccountsStoreActivity, nil, accounts)
}
