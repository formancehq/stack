package bankingcircle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/bankingcircle/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	accountsBalancesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts_and_balances"))...)
	accountsAttrs         = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))...)
	balancesAttrs         = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "balances"))...)
)

func taskFetchAccounts(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), accountsBalancesAttrs)
		}()

		for page := 1; ; page++ {
			pagedAccounts, err := client.GetAccounts(ctx, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsBalancesAttrs)
				return err
			}

			if len(pagedAccounts) == 0 {
				break
			}

			if err := ingestAccountsBatch(ctx, ingester, metricsRegistry, pagedAccounts); err != nil {
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
	ingester ingestion.Ingester,
	metricsRegistry metrics.MetricsRegistry,
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
				Reference: account.AccountID,
				Provider:  models.ConnectorProviderBankingCircle,
			},
			CreatedAt:    openingDate,
			Reference:    account.AccountID,
			Provider:     models.ConnectorProviderBankingCircle,
			DefaultAsset: models.Asset(account.Currency + "/2"),
			AccountName:  account.AccountDescription,
			Type:         models.AccountTypeInternal,
			RawData:      raw,
		})

		for _, balance := range account.Balances {
			var amount big.Float
			_, ok := amount.SetString(balance.IntraDayAmount.String())
			if !ok {
				return fmt.Errorf("failed to parse amount %s", balance.IntraDayAmount)
			}

			var amountInt big.Int
			amount.Mul(&amount, big.NewFloat(100)).Int(&amountInt)

			now := time.Now()
			balanceBatch = append(balanceBatch, &models.Balance{
				AccountID: models.AccountID{
					Reference: account.AccountID,
					Provider:  models.ConnectorProviderBankingCircle,
				},
				// Note(polo): same thing as payments
				// TODO(polo): do a complete pass on all connectors to
				// normalize the currencies
				Asset:         models.Asset(balance.Currency + "/2"),
				Balance:       &amountInt,
				CreatedAt:     now,
				LastUpdatedAt: now,
			})
		}
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountsBatch)), accountsAttrs)

	if err := ingester.IngestBalances(ctx, balanceBatch, false); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountsBatch)), balancesAttrs)

	return nil
}
