package stripe

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stripe/stripe-go/v72"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	externalAccountsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))...)
)

func FetchExternalAccountsTask(config Config, account string, client *client.DefaultClient) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), externalAccountsAttrs)
		}()

		tt := NewTimelineTrigger(
			logger,
			NewIngester(
				func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
					return nil

				},
				func(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error {
					return nil
				},
				func(ctx context.Context, batch []*stripe.ExternalAccount, commitState TimelineState, tail bool) error {
					if err := ingestExternalAccountsBatch(ctx, ingester, batch); err != nil {
						return err
					}
					metricsRegistry.ConnectorObjects().Add(ctx, int64(len(batch)), externalAccountsAttrs)
					return nil
				},
			),
			NewTimeline(client.ForAccount(account),
				config.TimelineConfig, task.MustResolveTo(ctx, resolver, TimelineState{})),
			TimelineTriggerTypeExternalAccounts,
		)

		if err := tt.Fetch(ctx); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, externalAccountsAttrs)
			return err
		}

		return nil
	}
}

func ingestExternalAccountsBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	accounts []*stripe.ExternalAccount,
) error {
	batch := ingestion.AccountBatch{}
	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		batch = append(batch, &models.Account{
			ID: models.AccountID{
				Reference: account.ID,
				Provider:  models.ConnectorProviderStripe,
			},
			CreatedAt:    time.Unix(account.BankAccount.Account.Created, 0),
			Reference:    account.ID,
			Provider:     models.ConnectorProviderStripe,
			DefaultAsset: currency.FormatAsset(string(account.BankAccount.Account.DefaultCurrency)),
			Type:         models.AccountTypeExternal,
			RawData:      raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
