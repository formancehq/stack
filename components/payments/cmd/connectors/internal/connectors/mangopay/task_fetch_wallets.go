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

		err := ingestWallets(ctx, client, config, userID, connectorID, scheduler, ingester)
		if err != nil {
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
) error {
	for currentPage := 1; ; currentPage++ {
		pagedWallets, err := client.GetWallets(ctx, userID, currentPage, pageSize)
		if err != nil {
			return err
		}

		if len(pagedWallets) == 0 {
			break
		}

		if err = handleWallets(
			ctx,
			config,
			userID,
			connectorID,
			ingester,
			scheduler,
			pagedWallets,
		); err != nil {
			return err
		}

		if len(pagedWallets) < pageSize {
			break
		}
	}

	return nil
}

func handleWallets(
	ctx context.Context,
	config *Config,
	userID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	pagedWallets []*client.Wallet,
) error {
	var accountBatch ingestion.AccountBatch
	var balanceBatch ingestion.BalanceBatch
	var transactionTasks []models.TaskDescriptor
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
			return fmt.Errorf("failed to parse amount: %s", wallet.Balance.Amount.String())
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
	}

	if err := ingester.IngestAccounts(ctx, accountBatch); err != nil {
		return err
	}

	if err := ingester.IngestBalances(ctx, balanceBatch, false); err != nil {
		return err
	}

	for _, transactionTask := range transactionTasks {
		err := scheduler.Schedule(ctx, transactionTask, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
			Duration:       config.PollingPeriod.Duration,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}
	}

	return nil
}
