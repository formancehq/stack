package bankingcircle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/bankingcircle/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchAccounts(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		for page := 1; ; page++ {
			pagedAccounts, err := client.GetAccounts(ctx, page)
			if err != nil {
				return err
			}

			if len(pagedAccounts) == 0 {
				break
			}

			if err := ingestAccountsBatch(ctx, ingester, pagedAccounts); err != nil {
				return err
			}
		}

		// We want to fetch payments after inserting all accounts in order to
		// ling them correctly
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

		openingDate, err := time.Parse("2006-01-02T15:04:05.999999999+00:00", account.OpeningDate)
		if err != nil {
			return fmt.Errorf("failed to parse opening date: %w", err)
		}

		batchElement := ingestion.AccountBatchElement{
			Account: &models.Account{
				ID: models.AccountID{
					Reference: account.AccountID,
					Provider:  models.ConnectorProviderBankingCircle,
				},
				CreatedAt:       openingDate,
				Reference:       account.AccountID,
				Provider:        models.ConnectorProviderBankingCircle,
				DefaultCurrency: account.Currency,
				AccountName:     account.AccountDescription,
				Type:            models.AccountTypeInternal,
				RawData:         raw,
			},
		}

		batch = append(batch, batchElement)
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
