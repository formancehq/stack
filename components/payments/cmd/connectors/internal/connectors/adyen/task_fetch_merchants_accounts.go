package adyen

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adyen/adyen-go-api-library/v7/src/management"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/adyen/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	pageSize = 100
)

var (
	accountsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "accounts"))...)
)

func taskFetchAccounts(client *client.Client) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchAccounts)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), accountsAttrs)
		}()

		for page := 1; ; page++ {
			pagedAccounts, err := client.GetMerchantAccounts(ctx, int32(page), pageSize)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
				return err
			}

			if err := ingestAccountsBatch(ctx, connectorID, ingester, pagedAccounts); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, accountsAttrs)
				return err
			}

			if len(pagedAccounts) < pageSize {
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
	accounts []management.Merchant,
) error {
	if len(accounts) == 0 {
		return nil
	}

	batch := ingestion.AccountBatch{}
	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return err
		}

		a := &models.Account{
			ID: models.AccountID{
				Reference:   *account.Id,
				ConnectorID: connectorID,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:   time.Now(),
			Reference:   *account.Id,
			ConnectorID: connectorID,
			Type:        models.AccountTypeInternal,
			RawData:     raw,
		}

		if account.Name != nil {
			a.AccountName = *account.Name
		}

		batch = append(batch, a)
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
