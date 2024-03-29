package generic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/generic/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

type fetchAccountsState struct {
	LastCreatedAt time.Time
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
			"generic.taskFetchAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchAccountsState{})

		newState, err := ingestAccounts(ctx, connectorID, client, ingester, scheduler, state)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingester.UpdateTaskState(ctx, newState); err != nil {
			otel.RecordError(span, err)
			return err
		}

		taskTransactions, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch transactions from client",
			Key:  taskNameFetchTransactions,
		})
		if err != nil {
			otel.RecordError(span, err)
			return errors.Wrap(task.ErrRetryable, err.Error())
		}

		err = scheduler.Schedule(ctx, taskTransactions, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			otel.RecordError(span, err)
			return errors.Wrap(task.ErrRetryable, err.Error())
		}

		return nil
	}
}

func ingestAccounts(
	ctx context.Context,
	connectorID models.ConnectorID,
	client *client.Client,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	state fetchAccountsState,
) (fetchAccountsState, error) {
	newState := fetchAccountsState{
		LastCreatedAt: state.LastCreatedAt,
	}

	balancesTasks := make([]models.TaskDescriptor, 0)
	for page := 1; ; page++ {
		accounts, err := client.ListAccounts(ctx, int64(page), pageSize, state.LastCreatedAt)
		if err != nil {
			return fetchAccountsState{}, err
		}

		if len(accounts) == 0 {
			break
		}

		accountsBatch := make([]*models.Account, 0, len(accounts))
		for _, account := range accounts {
			raw, err := json.Marshal(account)
			if err != nil {
				return fetchAccountsState{}, err
			}

			accountsBatch = append(accountsBatch, &models.Account{
				ID: models.AccountID{
					Reference:   account.Id,
					ConnectorID: connectorID,
				},
				ConnectorID: connectorID,
				CreatedAt:   account.CreatedAt,
				Reference:   account.Id,
				AccountName: account.AccountName,
				Type:        models.AccountTypeInternal,
				Metadata:    account.Metadata,
				RawData:     raw,
			})

			balanceTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch balances from client",
				Key:       taskNameFetchBalances,
				AccountID: account.Id,
			})
			if err != nil {
				return fetchAccountsState{}, err
			}

			balancesTasks = append(balancesTasks, balanceTask)

			newState.LastCreatedAt = account.CreatedAt
		}

		if err := ingester.IngestAccounts(ctx, ingestion.AccountBatch(accountsBatch)); err != nil {
			return fetchAccountsState{}, errors.Wrap(task.ErrRetryable, err.Error())
		}

		for _, balanceTask := range balancesTasks {
			if err := scheduler.Schedule(ctx, balanceTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			}); err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return fetchAccountsState{}, errors.Wrap(task.ErrRetryable, err.Error())
			}
		}
	}

	return newState, nil
}
