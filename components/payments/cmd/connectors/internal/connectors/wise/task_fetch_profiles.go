package wise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func taskFetchProfiles(wiseClient *client.Client) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		sp := trace.SpanFromContext(ctx)
		sp.SetName("wise.taskFetchProfiles")
		sp.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
		)

		if err := fetchProfiles(ctx, wiseClient, connectorID, ingester, scheduler); err != nil {
			otel.RecordError(sp, err)
			return err
		}

		return nil
	}
}

func fetchProfiles(
	ctx context.Context,
	wiseClient *client.Client,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
) error {
	profiles, err := wiseClient.GetProfiles(ctx)
	if err != nil {
		return err
	}

	var descriptors []models.TaskDescriptor
	for _, profile := range profiles {
		balances, err := wiseClient.GetBalances(ctx, profile.ID)
		if err != nil {
			return err
		}

		if err := ingestAccountsBatch(
			ctx,
			connectorID,
			ingester,
			profile.ID,
			balances,
		); err != nil {
			return err
		}

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

func ingestAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
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

		amount, err := currency.GetAmountWithPrecisionFromString(balance.Amount.Value.String(), precision)
		if err != nil {
			return err
		}

		now := time.Now()
		balancesBatch = append(balancesBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   fmt.Sprintf("%d", balance.ID),
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Amount.Currency),
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
