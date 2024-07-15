package atlar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/get-momo/atlar-v1-go-client/client/accounts"
	"github.com/get-momo/atlar-v1-go-client/client/external_accounts"
	"go.opentelemetry.io/otel/attribute"
)

func FetchAccountsTask(config Config, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"atlar.taskFetchAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		// Pagination works by cursor token.
		for token := ""; ; {
			requestCtx, cancel := contextutil.DetachedWithTimeout(ctx, 30*time.Second)
			defer cancel()
			pagedAccounts, err := client.GetV1Accounts(requestCtx, token, int64(config.PageSize))
			if err != nil {
				otel.RecordError(span, err)
				return err
			}

			token = pagedAccounts.Payload.NextToken

			if err := ingestAccountsBatch(ctx, connectorID, taskID, ingester, pagedAccounts, client); err != nil {
				otel.RecordError(span, err)
				return err
			}

			if token == "" {
				break
			}
		}

		// Pagination works by cursor token.
		for token := ""; ; {
			requestCtx, cancel := contextutil.DetachedWithTimeout(ctx, 30*time.Second)
			defer cancel()
			pagedExternalAccounts, err := client.GetV1ExternalAccounts(requestCtx, token, int64(config.PageSize))
			if err != nil {
				otel.RecordError(span, err)
				return err
			}

			token = pagedExternalAccounts.Payload.NextToken

			if err := ingestExternalAccountsBatch(ctx, connectorID, ingester, pagedExternalAccounts, client); err != nil {
				otel.RecordError(span, err)
				return err
			}

			if token == "" {
				break
			}
		}

		// Fetch payments after inserting all accounts in order to link them correctly
		taskPayments, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch payments from Atlar",
			Key:  taskNameFetchTransactions,
		})
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		err = scheduler.Schedule(ctx, taskPayments, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func ingestAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	taskID models.TaskID,
	ingester ingestion.Ingester,
	pagedAccounts *accounts.GetV1AccountsOK,
	client *client.Client,
) error {
	ctx, span := connectors.StartSpan(
		ctx,
		"atlar.taskFetchAccounts.ingestAccountsBatch",
		attribute.String("connectorID", connectorID.String()),
		attribute.String("taskID", taskID.String()),
	)
	defer span.End()

	accountsBatch := ingestion.AccountBatch{}
	balanceBatch := ingestion.BalanceBatch{}

	for _, account := range pagedAccounts.Payload.Items {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		createdAt, err := ParseAtlarTimestamp(account.Created)
		if err != nil {
			return fmt.Errorf("failed to parse opening date: %w", err)
		}

		requestCtx, cancel := contextutil.DetachedWithTimeout(ctx, 30*time.Second)
		defer cancel()
		thirdPartyResponse, err := client.GetV1BetaThirdPartiesID(requestCtx, account.ThirdPartyID)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference:   *account.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:    createdAt,
			Reference:    *account.ID,
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency),
			AccountName:  account.Name,
			Type:         models.AccountTypeInternal,
			Metadata:     ExtractAccountMetadata(account, thirdPartyResponse.Payload),
			RawData:      raw,
		})

		balance := account.Balance
		balanceTimestamp, err := ParseAtlarTimestamp(balance.Timestamp)
		if err != nil {
			return err
		}
		balanceBatch = append(balanceBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   *account.ID,
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, *balance.Amount.Currency),
			Balance:       big.NewInt(*balance.Amount.Value),
			CreatedAt:     balanceTimestamp,
			LastUpdatedAt: time.Now().UTC(),
			ConnectorID:   connectorID,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		return err
	}

	if err := ingester.IngestBalances(ctx, balanceBatch, false); err != nil {
		return err
	}

	return nil
}

func ingestExternalAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	pagedExternalAccounts *external_accounts.GetV1ExternalAccountsOK,
	client *client.Client,
) error {
	accountsBatch := ingestion.AccountBatch{}

	for _, externalAccount := range pagedExternalAccounts.Payload.Items {
		requestCtx, cancel := contextutil.DetachedWithTimeout(ctx, 30*time.Second)
		defer cancel()
		counterparty_response, err := client.GetV1CounterpartiesID(requestCtx, externalAccount.CounterpartyID)
		if err != nil {
			return err
		}
		counterparty := counterparty_response.Payload

		newAccount, err := ExternalAccountFromAtlarData(connectorID, externalAccount, counterparty)
		if err != nil {
			return err
		}

		accountsBatch = append(accountsBatch, newAccount)
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		return err
	}

	return nil
}
