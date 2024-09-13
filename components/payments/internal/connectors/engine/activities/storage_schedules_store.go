package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageSchedulesStore(ctx context.Context, schedule models.Schedule) error {
	return a.storage.SchedulesUpsert(ctx, schedule)
}

var StorageSchedulesStoreActivity = Activities{}.StorageSchedulesStore

func StorageSchedulesStore(ctx workflow.Context, schedule models.Schedule) error {
	return executeActivity(ctx, StorageSchedulesStoreActivity, nil, schedule)
}
