package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageBalancesStore(ctx context.Context, balances []models.Balance) error {
	return a.storage.BalancesUpsert(ctx, balances)
}

var StorageBalancesStoreActivity = Activities{}.StorageBalancesStore

func StorageBalancesStore(ctx workflow.Context, balances []models.Balance) error {
	return executeActivity(ctx, StorageBalancesStoreActivity, nil, balances)
}
