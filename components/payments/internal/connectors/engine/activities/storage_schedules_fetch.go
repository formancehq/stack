package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageSchedulesList(ctx context.Context, query storage.ListSchedulesQuery) (*bunpaginate.Cursor[models.Schedule], error) {
	return a.storage.SchedulesList(ctx, query)
}

var StorageSchedulesListActivity = Activities{}.StorageSchedulesList

func StorageSchedulesList(ctx workflow.Context, query storage.ListSchedulesQuery) (*bunpaginate.Cursor[models.Schedule], error) {
	ret := bunpaginate.Cursor[models.Schedule]{}
	if err := executeActivity(ctx, StorageSchedulesListActivity, &ret, query); err != nil {
		return nil, err
	}
	return &ret, nil
}
