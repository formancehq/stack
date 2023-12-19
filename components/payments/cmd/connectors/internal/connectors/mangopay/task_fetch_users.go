package mangopay

import (
	"context"
	"errors"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
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

			walletsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:   "Fetch wallets from client by user",
				Key:    taskNameFetchWallets,
				UserID: user.ID,
			})
			if err != nil {
				return err
			}

			err = scheduler.Schedule(ctx, walletsTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return err
			}

			bankAccountsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:   "Fetch bank accounts from client by user",
				Key:    taskNameFetchBankAccounts,
				UserID: user.ID,
			})
			if err != nil {
				return err
			}

			err = scheduler.Schedule(ctx, bankAccountsTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return err
			}
		}

		return nil
	}
}
