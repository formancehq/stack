package moneycorp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/connectors/moneycorp/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	recipientsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "recipients"))...)
)

func taskFetchRecipients(logger logging.Logger, client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchRecipients)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), recipientsAttrs)
		}()

		for page := 1; ; page++ {
			pagedRecipients, err := client.GetRecipients(ctx, accountID, page, pageSize)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, recipientsAttrs)
				return err
			}

			if len(pagedRecipients) == 0 {
				break
			}

			if err := ingestRecipientsBatch(ctx, ingester, pagedRecipients); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, recipientsAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedRecipients)), recipientsAttrs)

			if len(pagedRecipients) < pageSize {
				break
			}
		}

		return nil
	}
}

func ingestRecipientsBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	recipients []*client.Recipient,
) error {
	batch := ingestion.AccountBatch{}
	for _, recipient := range recipients {
		raw, err := json.Marshal(recipient)
		if err != nil {
			return err
		}

		createdAt, err := time.Parse("2006-01-02T15:04:05.999999999", recipient.Attributes.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to parse transaction date: %w", err)
		}

		batch = append(batch, &models.Account{
			ID: models.AccountID{
				Reference: recipient.ID,
				Provider:  models.ConnectorProviderMoneycorp,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:    createdAt,
			Reference:    recipient.ID,
			Provider:     models.ConnectorProviderMoneycorp,
			DefaultAsset: currency.FormatAsset(recipient.Attributes.BankAccountCurrency),
			AccountName:  recipient.Attributes.BankAccountName,
			Type:         models.AccountTypeExternal,
			RawData:      raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
