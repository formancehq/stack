package wise

import (
	"context"
	"errors"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

// taskMain is the main task of the connector. It launches the other tasks.
func taskMain(logger logging.Logger) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
	) error {
		logger.Info(taskNameMain)

		taskUsers, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch users from client",
			Key:  taskNameFetchProfiles,
		})
		if err != nil {
			return err
		}

		err = scheduler.Schedule(ctx, taskUsers, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			Restart:        true,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}

		return nil
	}
}
