package modulr

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

func taskFetchAccounts(config Config, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"modulr.taskFetchAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		err := fetchAccounts(
			ctx,
			config,
			client,
			connectorID,
			ingester,
			scheduler,
		)
		if err != nil {
			otel.RecordError(span, err)
			// Retry errors are handled by the function
			return err
		}

		return nil
	}
}

func fetchAccounts(
	ctx context.Context,
	config Config,
	client *client.Client,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
) error {
	for page := 0; ; page++ {
		pagedAccounts, err := client.GetAccounts(
			ctx,
			page,
			config.PageSize,
		)
		if err != nil {
			// Retry errors are handled by the client
			return err
		}

		if len(pagedAccounts.Content) == 0 {
			break
		}

		if err := ingestAccountsBatch(ctx, connectorID, ingester, pagedAccounts.Content); err != nil {
			return errors.Wrap(task.ErrRetryable, err.Error())
		}

		for _, account := range pagedAccounts.Content {
			transactionsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch transactions from client by account",
				Key:       taskNameFetchTransactions,
				AccountID: account.ID,
			})
			if err != nil {
				return errors.Wrap(task.ErrRetryable, err.Error())
			}

			err = scheduler.Schedule(ctx, transactionsTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return errors.Wrap(task.ErrRetryable, err.Error())
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

		openingDate, err := time.Parse(timeTemplate, account.CreatedDate)
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
		precision := supportedCurrenciesWithDecimal[account.Currency]

		amount, err := currency.GetAmountWithPrecisionFromString(account.Balance, precision)
		if err != nil {
			return fmt.Errorf("failed to parse amount %s: %w", account.Balance, err)
		}

		now := time.Now()
		balancesBatch = append(balancesBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   account.ID,
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency),
			Balance:       amount,
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
