package currencycloud

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchAccountsState struct {
	LastPage      int
	LastCreatedAt time.Time
}

func taskFetchAccounts(
	client *client.Client,
) task.Task {
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
			"currencycloud.taskFetchAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchAccountsState{})
		if state.LastPage == 0 {
			// First run, the first page for currencycloud starts at 1 and not 0
			state.LastPage = 1
		}

		newState, err := fetchAccount(ctx, client, connectorID, ingester, scheduler, state)
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

func fetchAccount(
	ctx context.Context,
	client *client.Client,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	state fetchAccountsState,
) (fetchAccountsState, error) {
	newState := fetchAccountsState{
		LastPage:      state.LastPage,
		LastCreatedAt: state.LastCreatedAt,
	}

	page := state.LastPage
	for {
		if page < 0 {
			break
		}

		pagedAccounts, nextPage, err := client.GetAccounts(ctx, page)
		if err != nil {
			return fetchAccountsState{}, err
		}

		page = nextPage

		batch := ingestion.AccountBatch{}
		for _, account := range pagedAccounts {
			switch account.CreatedAt.Compare(state.LastCreatedAt) {
			case -1, 0:
				// Account already ingested, skip
				continue
			default:
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
				CreatedAt:   account.CreatedAt,
				Reference:   account.ID,
				ConnectorID: connectorID,
				AccountName: account.AccountName,
				Type:        models.AccountTypeInternal,
				RawData:     raw,
			})

			newState.LastCreatedAt = account.CreatedAt
		}

		if err := ingester.IngestAccounts(ctx, batch); err != nil {
			return fetchAccountsState{}, err
		}
	}

	newState.LastPage = page

	taskTransactions, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name: "Fetch transactions from client",
		Key:  taskNameFetchTransactions,
	})
	if err != nil {
		return fetchAccountsState{}, err
	}

	err = scheduler.Schedule(ctx, taskTransactions, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
	})
	if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
		return fetchAccountsState{}, err
	}

	taskBalances, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name: "Fetch balances from client",
		Key:  taskNameFetchBalances,
	})
	if err != nil {
		return fetchAccountsState{}, err
	}

	err = scheduler.Schedule(ctx, taskBalances, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
	})
	if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
		return fetchAccountsState{}, err
	}

	return newState, nil
}
