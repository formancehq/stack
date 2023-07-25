package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/mangopay/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchWallets(logger logging.Logger, client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		logger.Info(taskNameFetchWallets)

		for page := 1; ; page++ {
			pagedWallets, err := client.GetWallets(ctx, userID, page)
			if err != nil {
				return err
			}

			if len(pagedWallets) == 0 {
				break
			}

			var accountBatch ingestion.AccountBatch
			var balanceBatch ingestion.BalanceBatch
			var transactionTasks []models.TaskDescriptor
			now := time.Now()
			for _, wallet := range pagedWallets {
				transactionTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
					Name:     "Fetch transactions from client by user and wallets",
					Key:      taskNameFetchTransactions,
					UserID:   userID,
					WalletID: wallet.ID,
				})
				if err != nil {
					return err
				}

				buf, err := json.Marshal(wallet)
				if err != nil {
					return err
				}

				transactionTasks = append(transactionTasks, transactionTask)
				accountBatch = append(accountBatch, &models.Account{
					ID: models.AccountID{
						Reference: wallet.ID,
						Provider:  models.ConnectorProviderMangopay,
					},
					CreatedAt:       time.Unix(wallet.CreationDate, 0),
					Reference:       wallet.ID,
					Provider:        models.ConnectorProviderMangopay,
					DefaultCurrency: wallet.Currency,
					AccountName:     wallet.Description,
					// Wallets are internal accounts on our side, since we
					// can have their balances.
					Type:    models.AccountTypeInternal,
					RawData: buf,
				})

				balanceBatch = append(balanceBatch, &models.Balance{
					AccountID: models.AccountID{
						Reference: wallet.ID,
						Provider:  models.ConnectorProviderMangopay,
					},
					Currency:      wallet.Balance.Currency,
					Balance:       wallet.Balance.Amount,
					CreatedAt:     now,
					LastUpdatedAt: now,
				})
			}

			if err := ingester.IngestAccounts(ctx, accountBatch); err != nil {
				return err
			}

			if err := ingester.IngestBalances(ctx, balanceBatch); err != nil {
				return err
			}

			for _, transactionTask := range transactionTasks {
				err = scheduler.Schedule(ctx, transactionTask, models.TaskSchedulerOptions{
					ScheduleOption: models.OPTIONS_RUN_NOW,
					Restart:        true,
				})
				if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
					return err
				}
			}
		}

		return nil
	}
}
