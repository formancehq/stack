package stripe

import (
	"context"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
)

func MainTask(logger logging.Logger) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
	) error {
		logger.Info("main task")

		taskAccounts, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch accounts from client",
			Key:  taskNameFetchAccounts,
		})
		if err != nil {
			return err
		}

		err = scheduler.Schedule(ctx, taskAccounts, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			Restart:        true,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}

		taskPayments, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch payments from client",
			Key:  taskNameFetchPayments,
		})
		if err != nil {
			return err
		}

		err = scheduler.Schedule(ctx, taskPayments, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			Restart:        true,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}

		return nil
	}
}
