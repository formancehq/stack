package wise

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	recipientAccountsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "recipient_accounts"))...)
)

func taskFetchRecipientAccounts(logger logging.Logger, client *client.Client, profileID uint64) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), recipientAccountsAttrs)
		}()

		recipientAccounts, err := client.GetRecipientAccounts(ctx, profileID)
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, recipientAccountsAttrs)
			return err
		}

		if err := ingestRecipientAccountsBatch(ctx, metricsRegistry, ingester, recipientAccounts); err != nil {
			return err
		}

		return nil
	}
}

func ingestRecipientAccountsBatch(
	ctx context.Context,
	metricsRegistry metrics.MetricsRegistry,
	ingester ingestion.Ingester,
	accounts []*client.RecipientAccount,
) error {
	accountsBatch := ingestion.AccountBatch{}
	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference: fmt.Sprintf("%d", account.ID),
				Provider:  models.ConnectorProviderWise,
			},
			CreatedAt:    time.Now(),
			Reference:    fmt.Sprintf("%d", account.ID),
			Provider:     models.ConnectorProviderWise,
			DefaultAsset: models.Asset(fmt.Sprintf("%s/2", account.Currency)),
			AccountName:  account.HolderName,
			Type:         models.AccountTypeExternal,
			RawData:      raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, recipientAccountsAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountsBatch)), recipientAccountsAttrs)

	return nil
}
