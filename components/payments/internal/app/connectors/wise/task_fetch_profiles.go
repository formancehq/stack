package wise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/wise/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	profilesAndBalancesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "profiles_and_balances"))...)
	profilesAttrs            = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "profiles"))...)
	balancesAttrs            = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "balances"))...)
)

func taskFetchProfiles(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), profilesAndBalancesAttrs)
		}()

		profiles, err := client.GetProfiles()
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, profilesAttrs)
			return err
		}

		var descriptors []models.TaskDescriptor
		for _, profile := range profiles {
			balances, err := client.GetBalances(ctx, profile.ID)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
				return err
			}

			if err := ingestAccountsBatch(ctx, metricsRegistry, ingester, balances); err != nil {
				return err
			}

			logger.Infof(fmt.Sprintf("scheduling fetch-transfers: %d", profile.ID))

			transferDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch transfers from client by profile",
				Key:       taskNameFetchTransfers,
				ProfileID: profile.ID,
			})
			if err != nil {
				return err
			}
			descriptors = append(descriptors, transferDescriptor)

			recipientAccountsDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:      "Fetch recipient accounts from client by profile",
				Key:       taskNameFetchRecipientAccounts,
				ProfileID: profile.ID,
			})
			if err != nil {
				return err
			}
			descriptors = append(descriptors, recipientAccountsDescriptor)
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
	metricsRegistry metrics.MetricsRegistry,
	ingester ingestion.Ingester,
	balances []*client.Balance,
) error {
	if len(balances) == 0 {
		return nil
	}

	accountsBatch := ingestion.AccountBatch{}
	balancesBatch := ingestion.BalanceBatch{}
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
			CreatedAt:    balance.CreationTime,
			Reference:    fmt.Sprintf("%d", balance.ID),
			Provider:     models.ConnectorProviderWise,
			DefaultAsset: models.Asset(fmt.Sprintf("%s/2", balance.Amount.Currency)),
			AccountName:  balance.Name,
			Type:         models.AccountTypeInternal,
			RawData:      raw,
		})

		var amount big.Float
		_, ok := amount.SetString(balance.Amount.Value.String())
		if !ok {
			return fmt.Errorf("failed to parse amount %s", balance.Amount.Value.String())
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(100)).Int(&amountInt)

		now := time.Now()
		balancesBatch = append(balancesBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference: fmt.Sprintf("%d", balance.ID),
				Provider:  models.ConnectorProviderWise,
			},
			Asset:         models.Asset(fmt.Sprintf("%s/2", balance.Amount.Currency)),
			Balance:       &amountInt,
			CreatedAt:     now,
			LastUpdatedAt: now,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, profilesAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountsBatch)), profilesAttrs)

	if err := ingester.IngestBalances(ctx, balancesBatch, false); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(balancesBatch)), balancesAttrs)

	return nil
}
