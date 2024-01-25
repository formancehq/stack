package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchWalletsState struct {
	LastPage         int       `json:"last_page"`
	LastCreationDate time.Time `json:"last_creation_date"`
}

func taskFetchWallets(client *client.Client, config *Config, userID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"mangopay.taskFetchWallets",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("userID", userID),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchWalletsState{})
		if state.LastPage == 0 {
			// If last page is 0, it means we haven't fetched any wallets yet.
			// Mangopay pages starts at 1.
			state.LastPage = 1
		}

		newState, err := ingestWallets(ctx, client, config, userID, connectorID, scheduler, ingester, state)
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

func ingestWallets(
	ctx context.Context,
	client *client.Client,
	config *Config,
	userID string,
	connectorID models.ConnectorID,
	scheduler task.Scheduler,
	ingester ingestion.Ingester,
	state fetchWalletsState,
) (fetchWalletsState, error) {
	var currentPage int

	newState := fetchWalletsState{
		LastCreationDate: state.LastCreationDate,
	}

	for currentPage = state.LastPage; ; currentPage++ {
		pagedWallets, err := client.GetWallets(ctx, userID, currentPage, pageSize)
		if err != nil {
			return fetchWalletsState{}, err
		}

		if len(pagedWallets) == 0 {
			break
		}

		lastCreationDate, err := handleWallets(
			ctx,
			config,
			userID,
			connectorID,
			ingester,
			scheduler,
			pagedWallets,
			state,
		)
		if err != nil {
			return fetchWalletsState{}, err
		}
		newState.LastCreationDate = lastCreationDate

		if len(pagedWallets) < pageSize {
			break
		}
	}

	newState.LastPage = currentPage

	return newState, nil
}

func handleWallets(
	ctx context.Context,
	config *Config,
	userID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	pagedWallets []*client.Wallet,
	state fetchWalletsState,
) (time.Time, error) {
	var accountBatch ingestion.AccountBatch
	var balanceBatch ingestion.BalanceBatch
	var transactionTasks []models.TaskDescriptor
	var lastCreationDate time.Time
	for _, wallet := range pagedWallets {
		creationDate := time.Unix(wallet.CreationDate, 0)
		switch creationDate.Compare(state.LastCreationDate) {
		case -1, 0:
			// creationDate <= state.LastCreationDate, nothing to do,
			// we already processed wallets.
			continue
		default:
		}

		transactionTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:     "Fetch transactions from client by user and wallets",
			Key:      taskNameFetchTransactions,
			UserID:   userID,
			WalletID: wallet.ID,
		})
		if err != nil {
			return time.Time{}, err
		}

		buf, err := json.Marshal(wallet)
		if err != nil {
			return time.Time{}, err
		}

		transactionTasks = append(transactionTasks, transactionTask)
		accountBatch = append(accountBatch, &models.Account{
			ID: models.AccountID{
				Reference:   wallet.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:    time.Unix(wallet.CreationDate, 0),
			Reference:    wallet.ID,
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, wallet.Currency),
			AccountName:  wallet.Description,
			// Wallets are internal accounts on our side, since we
			// can have their balances.
			Type: models.AccountTypeInternal,
			Metadata: map[string]string{
				"user_id": userID,
			},
			RawData: buf,
		})

		var amount big.Int
		_, ok := amount.SetString(wallet.Balance.Amount.String(), 10)
		if !ok {
			return time.Time{}, fmt.Errorf("failed to parse amount: %s", wallet.Balance.Amount.String())
		}

		now := time.Now()
		balanceBatch = append(balanceBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   wallet.ID,
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, wallet.Balance.Currency),
			Balance:       &amount,
			CreatedAt:     now,
			LastUpdatedAt: now,
			ConnectorID:   connectorID,
		})

		lastCreationDate = creationDate
	}

	if err := ingester.IngestAccounts(ctx, accountBatch); err != nil {
		return time.Time{}, err
	}

	if err := ingester.IngestBalances(ctx, balanceBatch, false); err != nil {
		return time.Time{}, err
	}

	for _, transactionTask := range transactionTasks {
		err := scheduler.Schedule(ctx, transactionTask, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
			Duration:       config.PollingPeriod.Duration,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return time.Time{}, err
		}
	}

	return lastCreationDate, nil
}
