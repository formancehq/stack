package modulr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/internal/models"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchAccounts(logger logging.Logger, config Config, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		logger.Info(taskNameFetchAccounts)

		for page := 0; ; page++ {
			pagedAccounts, err := client.GetAccounts(ctx, page, config.PageSize)
			if err != nil {
				return err
			}

			if len(pagedAccounts.Content) == 0 {
				break
			}

			if err := ingestAccountsBatch(ctx, connectorID, ingester, pagedAccounts.Content); err != nil {
				return err
			}

			for _, account := range pagedAccounts.Content {
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
					RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
				})
				if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
					return err
				}
			}

			if len(pagedAccounts.Content) < config.PageSize {
				break
			}

			if page+1 >= pagedAccounts.TotalPages {
				// Modulr paging starts at 0, so the last page is TotalPages - 1.
				break
			}
		}

		return nil
	}
}

func ingestAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	accounts []*client.Account,
) error {
	accountsBatch := ingestion.AccountBatch{}
	balancesBatch := ingestion.BalanceBatch{}

	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		openingDate, err := time.Parse("2006-01-02T15:04:05.999999999+0000", account.CreatedDate)
		if err != nil {
			return err
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference:   account.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:    openingDate,
			Reference:    account.ID,
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency),
			AccountName:  account.Name,
			Type:         models.AccountTypeInternal,
			RawData:      raw,
		})

		// No need to check if the currency is supported for accounts and
		// balances.
		precision, _ := supportedCurrenciesWithDecimal[account.Currency]

		var amount big.Float
		_, ok := amount.SetString(account.Balance)
		if !ok {
			return fmt.Errorf("failed to parse amount %s", account.Balance)
		}

		var balance big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&balance)

		now := time.Now()
		balancesBatch = append(balancesBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   account.ID,
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency),
			Balance:       &balance,
			CreatedAt:     now,
			LastUpdatedAt: now,
			ConnectorID:   connectorID,
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
