package atlar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	atlar_client "github.com/get-momo/atlar-v1-go-client/client"
	"github.com/get-momo/atlar-v1-go-client/client/accounts"
	"github.com/get-momo/atlar-v1-go-client/client/counterparties"
	"github.com/get-momo/atlar-v1-go-client/client/external_accounts"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	accountsBalancesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts_and_balances"))...)
	accountsAttrs         = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))...)
	balancesAttrs         = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "balances"))...)
)

func FetchAccountsTask(config Config, client *atlar_client.Rest) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), accountsBalancesAttrs)
		}()

		// Pagination works by cursor token.
		accountsParams := accounts.GetV1AccountsParams{
			Context: ctx,
			Limit:   pointer.For(int64(config.ApiConfig.PageSize)),
		}
		for token := ""; ; {
			limit := int64(config.PageSize)
			accountsParams.Token = &token
			accountsParams.Limit = &limit
			pagedAccounts, err := client.Accounts.GetV1Accounts(&accountsParams)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsBalancesAttrs)
				return err
			}

			token = pagedAccounts.Payload.NextToken

			if err := ingestAccountsBatch(ctx, connectorID, ingester, metricsRegistry, pagedAccounts); err != nil {
				return err
			}

			if token == "" {
				break
			}
		}

		// Pagination works by cursor token.
		externalAccountsParams := external_accounts.GetV1ExternalAccountsParams{
			Context: ctx,
			Limit:   pointer.For(int64(config.ApiConfig.PageSize)),
		}
		for token := ""; ; {
			limit := int64(config.PageSize)
			accountsParams.Token = &token
			accountsParams.Limit = &limit
			pagedExternalAccounts, err := client.ExternalAccounts.GetV1ExternalAccounts(&externalAccountsParams)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsBalancesAttrs)
				return err
			}

			token = pagedExternalAccounts.Payload.NextToken

			if err := ingestExternalAccountsBatch(ctx, connectorID, ingester, metricsRegistry, pagedExternalAccounts, client); err != nil {
				return err
			}

			if token == "" {
				break
			}
		}

		// Fetch payments after inserting all accounts in order to link them correctly
		taskPayments, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch payments from Atlar",
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
	metricsRegistry metrics.MetricsRegistry,
	pagedAccounts *accounts.GetV1AccountsOK,
) error {
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
			Metadata:     ExtractAccountMetadata(account),
			RawData:      raw,
		})

		balance := account.Balance
		// No need to check if the currency is supported for accounts and
		// balances.
		precision := supportedCurrenciesWithDecimal[*balance.Amount.Currency]

		var amount big.Float
		_, ok := amount.SetString(*balance.Amount.StringValue)
		if !ok {
			return fmt.Errorf("failed to parse amount %s", *balance.Amount.StringValue)
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

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
			Balance:       &amountInt,
			CreatedAt:     balanceTimestamp,
			LastUpdatedAt: time.Now().UTC(),
			ConnectorID:   connectorID,
		})
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

func ExtractAccountMetadata(account *atlar_models.Account) metadata.Metadata {
	result := metadata.Metadata{}
	result = result.Merge(ComputeAccountMetadataBool("fictive", account.Fictive))
	result = result.Merge(ComputeAccountMetadata("bank/id", account.Bank.ID))
	result = result.Merge(ComputeAccountMetadata("bank/name", account.Bank.Name))
	result = result.Merge(ComputeAccountMetadata("bank/bic", account.Bank.Bic))
	result = result.Merge(IdentifiersToMetadata(account.Identifiers))
	result = result.Merge(ComputeAccountMetadata("alias", account.Alias))
	result = result.Merge(ComputeAccountMetadata("owner/name", account.Owner.Name))
	return result
}

func IdentifiersToMetadata(identifiers []*atlar_models.AccountIdentifier) metadata.Metadata {
	result := metadata.Metadata{}
	for _, i := range identifiers {
		result = result.Merge(ComputeAccountMetadata(
			fmt.Sprintf("identifier/%s/%s", *i.Market, *i.Type),
			*i.Number,
		))
	}
	return result
}

type AtlarExternalAccountAndCounterparty struct {
	ExternalAccount atlar_models.ExternalAccount `json:"externalAccount" yaml:"externalAccount" bson:"externalAccount"`
	Counterparty    atlar_models.Counterparty    `json:"counterparty" yaml:"counterparty" bson:"counterparty"`
}

func ingestExternalAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	metricsRegistry metrics.MetricsRegistry,
	pagedExternalAccounts *external_accounts.GetV1ExternalAccountsOK,
	client *atlar_client.Rest,
) error {
	accountsBatch := ingestion.AccountBatch{}

	for _, externalAccount := range pagedExternalAccounts.Payload.Items {
		counterpartyParams := counterparties.GetV1CounterpartiesIDParams{
			Context: ctx,
			ID:      externalAccount.CounterpartyID,
		}
		counterparty_response, err := client.Counterparties.GetV1CounterpartiesID(&counterpartyParams)
		if err != nil {
			return err
		}
		counterparty := counterparty_response.Payload

		raw, err := json.Marshal(AtlarExternalAccountAndCounterparty{ExternalAccount: *externalAccount, Counterparty: *counterparty})
		if err != nil {
			return err
		}

		createdAt, err := ParseAtlarTimestamp(externalAccount.Created)
		if err != nil {
			return fmt.Errorf("failed to parse opening date: %w", err)
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference:   externalAccount.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:   createdAt,
			Reference:   externalAccount.ID,
			ConnectorID: connectorID,
			// DefaultAsset: left empty because the information is not provided by Atlar,
			AccountName: counterparty.Name, // TODO: is that okay? External accounts do not have a name at Atlar.
			Type:        models.AccountTypeExternal,
			Metadata:    ExtractExternalAccountAndCounterpartyMetadata(externalAccount, counterparty),
			RawData:     raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountsBatch)), accountsAttrs)

	return nil
}

func ExtractExternalAccountAndCounterpartyMetadata(externalAccount *atlar_models.ExternalAccount, counterparty *atlar_models.Counterparty) metadata.Metadata {
	result := metadata.Metadata{}
	result = result.Merge(ComputeAccountMetadata("bank/id", externalAccount.Bank.ID))
	result = result.Merge(ComputeAccountMetadata("bank/name", externalAccount.Bank.Name))
	result = result.Merge(ComputeAccountMetadata("bank/bic", externalAccount.Bank.Bic))
	result = result.Merge(IdentifiersToMetadata(externalAccount.Identifiers))
	result = result.Merge(ComputeAccountMetadata("owner/name", counterparty.Name))
	result = result.Merge(ComputeAccountMetadata("owner/type", counterparty.PartyType))
	result = result.Merge(ComputeAccountMetadata("owner/contact/email", counterparty.ContactDetails.Email))
	result = result.Merge(ComputeAccountMetadata("owner/contact/phone", counterparty.ContactDetails.Phone))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/streetName", counterparty.ContactDetails.Address.StreetName))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/streetNumber", counterparty.ContactDetails.Address.StreetNumber))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/city", counterparty.ContactDetails.Address.City))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/postalCode", counterparty.ContactDetails.Address.PostalCode))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/country", counterparty.ContactDetails.Address.Country))
	return result
}
