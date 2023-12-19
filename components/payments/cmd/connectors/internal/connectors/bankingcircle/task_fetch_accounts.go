package bankingcircle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchAccounts(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
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

			if err := ingestAccountsBatch(ctx, connectorID, ingester, pagedAccounts); err != nil {
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
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
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
	balanceBatch := ingestion.BalanceBatch{}

	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		openingDate, err := time.Parse("2006-01-02T15:04:05.999999999+00:00", account.OpeningDate)
		if err != nil {
			return fmt.Errorf("failed to parse opening date: %w", err)
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference:   account.AccountID,
				ConnectorID: connectorID,
			},
			CreatedAt:    openingDate,
			Reference:    account.AccountID,
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency),
			AccountName:  account.AccountDescription,
			Type:         models.AccountTypeInternal,
			RawData:      raw,
		})

		for _, balance := range account.Balances {
			// No need to check if the currency is supported for accounts and
			// balances.
			precision := supportedCurrenciesWithDecimal[balance.Currency]

			var amount big.Float
			_, ok := amount.SetString(balance.IntraDayAmount.String())
			if !ok {
				return fmt.Errorf("failed to parse amount %s", balance.IntraDayAmount)
			}

			var amountInt big.Int
			amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

			now := time.Now()
			balanceBatch = append(balanceBatch, &models.Balance{
				AccountID: models.AccountID{
					Reference:   account.AccountID,
					ConnectorID: connectorID,
				},
				Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Currency),
				Balance:       &amountInt,
				CreatedAt:     now,
				LastUpdatedAt: now,
				ConnectorID:   connectorID,
			})
		}
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		return err
	}

	if err := ingester.IngestBalances(ctx, balanceBatch, false); err != nil {
		return err
	}

	return nil
}
