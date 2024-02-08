package modulr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchAccountsState struct {
	LastCreatedTime time.Time `json:"last_created_date"`
}

func (s *fetchAccountsState) UpdateLatest(latest *client.Account) error {
	createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", latest.CreatedDate)
	if err != nil {
		return err
	}
	if createdTime.After(s.LastCreatedTime) {
		s.LastCreatedTime = createdTime
	}
	return nil
}

func (s *fetchAccountsState) IsNew(account *client.Account) (bool, error) {
	createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", account.CreatedDate)
	if err != nil {
		return false, err
	}
	return createdTime.After(s.LastCreatedTime), nil
}

func (s *fetchAccountsState) FilterNew(accounts []*client.Account) ([]*client.Account, error) {
	// accounts are assumed to be sorted by creation date.
	firstNewIdx := len(accounts)
	for idx, account := range accounts {
		isNew, err := s.IsNew(account)
		if err != nil {
			return nil, err
		}
		if isNew {
			firstNewIdx = idx
			break
		}
	}
	return accounts[firstNewIdx:], nil
}

func (s *fetchAccountsState) GetFilterValue() string {
	if s.LastCreatedTime.IsZero() {
		return ""
	}
	return s.LastCreatedTime.Format("2006-01-02T15:04:05-0700")
}

func taskFetchAccounts(config Config, client *client.Client) task.Task {
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
			"modulr.taskFetchAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state, err := fetchAccounts(
			ctx,
			config,
			client,
			connectorID,
			ingester,
			scheduler,
			task.MustResolveTo(ctx, resolver, fetchAccountsState{}),
		)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingester.UpdateTaskState(ctx, state); err != nil {
			otel.RecordError(span, err)
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
	state fetchAccountsState,
) (fetchAccountsState, error) {
	newState := state
	for page := 0; ; page++ {
		pagedAccounts, err := client.GetAccounts(
			ctx,
			page,
			config.PageSize,
			state.GetFilterValue(),
		)
		if err != nil {
			return newState, err
		}
		accounts, err := state.FilterNew(pagedAccounts.Content)
		if err != nil {
			return newState, err
		}
		if len(accounts) == 0 {
			break
		}
		if err := ingestAccountsBatch(ctx, connectorID, ingester, accounts); err != nil {
			return newState, err
		}

		for _, account := range accounts {
			if err := newState.UpdateLatest(account); err != nil {
				return newState, err
			}
			transactionsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch transactions from client by account",
				Key:       taskNameFetchTransactions,
				AccountID: account.ID,
			})
			if err != nil {
				return newState, err
			}

			err = scheduler.Schedule(ctx, transactionsTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
				Duration:       config.PollingPeriod.Duration,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return newState, err
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

	return newState, nil
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
