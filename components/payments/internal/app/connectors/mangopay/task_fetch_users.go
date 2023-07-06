package mangopay

import (
	"context"
	"errors"

	"github.com/formancehq/payments/internal/app/connectors/mangopay/client"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchUsers(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
	) error {
		logger.Info(taskNameFetchUsers)

		users, err := client.GetAllUsers(ctx)
		if err != nil {
			return err
		}

		for _, user := range users {
			logger.Infof("scheduling fetch-transactions: %s", user.ID)

			transactionsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:   "Fetch transactions from client by user",
				Key:    taskNameFetchTransactions,
				UserID: user.ID,
			})
			if err != nil {
				return err
			}

			err = scheduler.Schedule(ctx, transactionsTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				Restart:        true,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return err
			}
		}

		return nil
	}
}
