package wise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	profilesAndBalancesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "profiles_and_balances"))...)
	profilesAttrs            = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "profiles"))...)
	balancesAttrs            = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "balances"))...)
)

func taskFetchProfiles(wiseClient *client.Client) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), profilesAndBalancesAttrs)
		}()

		profiles, err := wiseClient.GetProfiles()
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, profilesAttrs)
			return err
		}

		var descriptors []models.TaskDescriptor
		for _, profile := range profiles {
			balances, err := wiseClient.GetBalances(ctx, profile.ID)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, balancesAttrs)
				return err
			}

			if err := ingestAccountsBatch(
				ctx,
				logger,
				connectorID,
				metricsRegistry,
				ingester,
				profile.ID,
				balances,
			); err != nil {
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
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
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
	logger logging.Logger,
	connectorID models.ConnectorID,
	metricsRegistry metrics.MetricsRegistry,
	ingester ingestion.Ingester,
	profileID uint64,
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

		precision, ok := supportedCurrenciesWithDecimal[balance.Amount.Currency]
		if !ok {
			logger.Errorf("currency %s is not supported", balance.Amount.Currency)
			metricsRegistry.ConnectorCurrencyNotSupported().Add(ctx, 1, metric.WithAttributes(connectorAttrs...))
			continue
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference:   fmt.Sprintf("%d", balance.ID),
				ConnectorID: connectorID,
			},
			CreatedAt:    balance.CreationTime,
			Reference:    fmt.Sprintf("%d", balance.ID),
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Amount.Currency),
			AccountName:  balance.Name,
			Type:         models.AccountTypeInternal,
			Metadata: map[string]string{
				"profile_id": strconv.FormatUint(profileID, 10),
			},
			RawData: raw,
		})

		var amount big.Float
		_, ok = amount.SetString(balance.Amount.Value.String())
		if !ok {
			return fmt.Errorf("failed to parse amount %s", balance.Amount.Value.String())
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

		now := time.Now()
		balancesBatch = append(balancesBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   fmt.Sprintf("%d", balance.ID),
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Amount.Currency),
			Balance:       &amountInt,
			CreatedAt:     now,
			LastUpdatedAt: now,
			ConnectorID:   connectorID,
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
