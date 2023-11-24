package atlar

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client/accounts"
	atlar_models "github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/models"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	accountsBalancesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts_and_balances"))...)
	accountsAttrs         = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))...)
	balancesAttrs         = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "balances"))...)
)

func FetchAccountsTask(config Config, client *client.Client) task.Task {
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

		params := accounts.GetV1AccountsParams{
			Context: ctx,
		}
		for token := ""; ; {
			params.Token = &token
			pagedAccounts, err := client.Accounts.GetV1Accounts(&params)
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

		openingDate, err := time.Parse("2006-01-02T15:04:05.999999999Z", account.Created)
		if err != nil {
			return fmt.Errorf("failed to parse opening date: %w", err)
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference:   *account.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:    openingDate,
			Reference:    *filterIbanFromIdentifiers(account.Identifiers).Number,
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency),
			AccountName:  account.Name,
			Type:         models.AccountTypeExternal,
			Metadata:     map[string]string{},
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

		now := time.Now()
		balanceBatch = append(balanceBatch, &models.Balance{
			AccountID: models.AccountID{
				Reference:   *account.ID,
				ConnectorID: connectorID,
			},
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, *balance.Amount.Currency),
			Balance:       &amountInt,
			CreatedAt:     now,
			LastUpdatedAt: now,
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

func filterIbanFromIdentifiers(identifiers []*atlar_models.AccountIdentifier) *atlar_models.AccountIdentifier {
	for _, i := range identifiers {
		if *i.Type == atlar_models.AccountIdentifierTypeIBAN {
			return i
		}
	}
	return nil
}
