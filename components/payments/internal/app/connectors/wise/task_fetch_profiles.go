package wise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/wise/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchProfiles(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		profiles, err := client.GetProfiles()
		if err != nil {
			return err
		}

		var descriptors []models.TaskDescriptor
		for _, profile := range profiles {
			balances, err := client.GetBalances(ctx, profile.ID)
			if err != nil {
				return err
			}

			if err := ingestAccountsBatch(ctx, ingester, balances); err != nil {
				return err
			}

			logger.Infof(fmt.Sprintf("scheduling fetch-transfers: %d", profile.ID))

			descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch transfers from client by profile",
				Key:       taskNameFetchTransfers,
				ProfileID: profile.ID,
			})
			if err != nil {
				return err
			}

			descriptors = append(descriptors, descriptor)
		}

		for _, descriptor := range descriptors {
			err = scheduler.Schedule(ctx, descriptor, models.TaskSchedulerOptions{
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

func ingestAccountsBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	balances []*client.Balance,
) error {
	if len(balances) == 0 {
		return nil
	}

	accountsBatch := ingestion.AccountBatch{}
	balancesBatch := ingestion.BalanceBatch{}
	now := time.Now()
	for _, balance := range balances {
		raw, err := json.Marshal(balance)
		if err != nil {
			return err
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference: fmt.Sprintf("%d", balance.ID),
				Provider:  models.ConnectorProviderWise,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:       balance.CreationTime,
			Reference:       fmt.Sprintf("%d", balance.ID),
			Provider:        models.ConnectorProviderWise,
			DefaultCurrency: balance.Currency,
			AccountName:     balance.Name,
			Type:            models.AccountTypeInternal,
			RawData:         raw,
		})

		balancesBatch = append(balancesBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference: fmt.Sprintf("%d", balance.ID),
				Provider:  models.ConnectorProviderWise,
			},
			Currency:      models.PaymentAsset(fmt.Sprintf("%s/2", balance.Amount.Currency)).String(),
			Balance:       int64(balance.Amount.Value * 100),
			CreatedAt:     now,
			LastUpdatedAt: now,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		return err
	}

	if err := ingester.IngestBalances(ctx, balancesBatch, false); err != nil {
		return err
	}

	return nil
}
