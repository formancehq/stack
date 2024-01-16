package mangopay

import (
	"context"
	"errors"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func taskFetchUsers(client *client.Client) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("mangopay.taskFetchUsers")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
		)

		if err := fetchUsers(ctx, client, connectorID, scheduler); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchUsers(
	ctx context.Context,
	client *client.Client,
	connectorID models.ConnectorID,
	scheduler task.Scheduler,
) error {
	users, err := client.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
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
