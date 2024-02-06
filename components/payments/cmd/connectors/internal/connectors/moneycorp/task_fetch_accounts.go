package moneycorp

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchAccountsState struct {
	LastPage int `json:"last_page"`
	// Moneycorp does not send the creation date for accounts, but we can still
	// sort by ID created (which is incremental when creating accounts).
	LastIDCreated string `json:"last_id_created"`
}

func taskFetchAccounts(client *client.Client, config *Config) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"moneycorp.taskFetchAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchAccountsState{})

		newState, err := fetchAccounts(ctx, config, client, connectorID, ingester, scheduler, state)
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

func fetchAccounts(
	ctx context.Context,
	config *Config,
	client *client.Client,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	state fetchAccountsState,
) (fetchAccountsState, error) {
	newState := fetchAccountsState{
		LastPage:      state.LastPage,
		LastIDCreated: state.LastIDCreated,
	}

	for page := state.LastPage; ; page++ {
		newState.LastPage = page

		pagedAccounts, err := client.GetAccounts(ctx, page, pageSize)
		if err != nil {
			return fetchAccountsState{}, err
		}

		if len(pagedAccounts) == 0 {
			break
		}

		batch := ingestion.AccountBatch{}
		transactionTasks := []models.TaskDescriptor{}
		balanceTasks := []models.TaskDescriptor{}
		recipientTasks := []models.TaskDescriptor{}
		for _, account := range pagedAccounts {
			if account.ID <= state.LastIDCreated {
				continue
			}

			raw, err := json.Marshal(account)
			if err != nil {
				return fetchAccountsState{}, err
			}

			batch = append(batch, &models.Account{
				ID: models.AccountID{
					Reference:   account.ID,
					ConnectorID: connectorID,
				},
				// Moneycorp does not send the opening date of the account
				CreatedAt:   time.Now().UTC(),
				Reference:   account.ID,
				ConnectorID: connectorID,
				AccountName: account.Attributes.AccountName,
				Type:        models.AccountTypeInternal,
				RawData:     raw,
			})

			transactionTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch transactions from client by account",
				Key:       taskNameFetchTransactions,
				AccountID: account.ID,
			})
			if err != nil {
				return fetchAccountsState{}, err
			}
			transactionTasks = append(transactionTasks, transactionTask)

			balanceTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch balances from client by account",
				Key:       taskNameFetchBalances,
				AccountID: account.ID,
			})
			if err != nil {
				return fetchAccountsState{}, err
			}
			balanceTasks = append(balanceTasks, balanceTask)

			recipientTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch recipients from client",
				Key:       taskNameFetchRecipients,
				AccountID: account.ID,
			})
			if err != nil {
				return fetchAccountsState{}, err
			}
			recipientTasks = append(recipientTasks, recipientTask)

			newState.LastIDCreated = account.ID
		}

		if err := ingester.IngestAccounts(ctx, batch); err != nil {
			return fetchAccountsState{}, err
		}

		for _, transactionTask := range transactionTasks {
			if err := scheduler.Schedule(ctx, transactionTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
				Duration:       config.PollingPeriod.Duration,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			}); err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return fetchAccountsState{}, err
			}
		}

		for _, balanceTask := range balanceTasks {
			if err := scheduler.Schedule(ctx, balanceTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
				Duration:       config.PollingPeriod.Duration,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			}); err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return fetchAccountsState{}, err
			}
		}

		for _, recipientTask := range recipientTasks {
			if err := scheduler.Schedule(ctx, recipientTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
				Duration:       config.PollingPeriod.Duration,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			}); err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return fetchAccountsState{}, err
			}
		}

		if len(pagedAccounts) < pageSize {
			break
		}
	}

	return newState, nil
}
