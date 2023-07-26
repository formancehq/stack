package moneycorp

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/moneycorp/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchAccounts(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		logger.Info(taskNameFetchAccounts)

		for page := 1; ; page++ {
			pagedAccounts, err := client.GetAccounts(ctx, page, pageSize)
			if err != nil {
				return err
			}

			if len(pagedAccounts) == 0 {
				break
			}

			if err := ingestAccountsBatch(ctx, ingester, pagedAccounts); err != nil {
				return err
			}

			for _, account := range pagedAccounts {
				logger.Infof("scheduling fetch-transactions: %s", account.ID)

				transactionsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
					Name:      "Fetch transactions from client by account",
					Key:       taskNameFetchTransactions,
					AccountID: account.ID,
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

				balancesTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
					Name:      "Fetch balances from client by account",
					Key:       taskNameFetchBalances,
					AccountID: account.ID,
				})

				if err != nil {
					return err
				}
				err = scheduler.Schedule(ctx, balancesTask, models.TaskSchedulerOptions{
					ScheduleOption: models.OPTIONS_RUN_NOW,
					Restart:        true,
				})
				if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
					return err
				}
			}

			if len(pagedAccounts) < pageSize {
				break
			}
		}

		return nil
	}
}

func ingestAccountsBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	accounts []*client.Account,
) error {
	batch := ingestion.AccountBatch{}
	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		batch = append(batch, &models.Account{
			ID: models.AccountID{
				Reference: account.ID,
				Provider:  models.ConnectorProviderMoneycorp,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:   time.Now(),
			Reference:   account.ID,
			Provider:    models.ConnectorProviderMoneycorp,
			AccountName: account.Attributes.AccountName,
			Type:        models.AccountTypeInternal,
			RawData:     raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
