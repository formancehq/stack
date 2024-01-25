package mangopay

import (
	"context"
	"errors"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchUsersState struct {
	LastPage         int       `json:"last_page"`
	LastCreationDate time.Time `json:"last_creation_date"`
}

func taskFetchUsers(client *client.Client, config *Config) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"mangopay.taskFetchUsers",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchUsersState{})
		if state.LastPage == 0 {
			// If last page is 0, it means we haven't fetched any users yet.
			// Mangopay pages starts at 1.
			state.LastPage = 1
		}

		newState, err := ingestUsers(ctx, client, config, connectorID, scheduler, state)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingester.UpdateTaskState(ctx, newState); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func ingestUsers(
	ctx context.Context,
	client *client.Client,
	config *Config,
	connectorID models.ConnectorID,
	scheduler task.Scheduler,
	state fetchUsersState,
) (fetchUsersState, error) {
	users, lastPage, err := client.GetAllUsers(ctx, state.LastPage, pageSize)
	if err != nil {
		return fetchUsersState{}, err
	}

	newState := fetchUsersState{
		LastPage:         lastPage,
		LastCreationDate: state.LastCreationDate,
	}

	for _, user := range users {
		userCreationDate := time.Unix(user.CreationDate, 0)
		switch userCreationDate.Compare(state.LastCreationDate) {
		case -1, 0:
			// creationDate <= state.LastCreationDate, nothing to do,
			// we already processed this user.
			continue
		default:
		}

		walletsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:   "Fetch wallets from client by user",
			Key:    taskNameFetchWallets,
			UserID: user.ID,
		})
		if err != nil {
			return fetchUsersState{}, err
		}

		err = scheduler.Schedule(ctx, walletsTask, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
			Duration:       config.PollingPeriod.Duration,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return fetchUsersState{}, err
		}

		bankAccountsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:   "Fetch bank accounts from client by user",
			Key:    taskNameFetchBankAccounts,
			UserID: user.ID,
		})
		if err != nil {
			return fetchUsersState{}, err
		}

		err = scheduler.Schedule(ctx, bankAccountsTask, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
			Duration:       config.PollingPeriod.Duration,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return fetchUsersState{}, err
		}

		newState.LastCreationDate = userCreationDate
	}

	return newState, nil
}
