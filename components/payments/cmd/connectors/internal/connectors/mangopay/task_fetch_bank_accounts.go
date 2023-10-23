package mangopay

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	bankAccountsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "bank_accounts"))...)
)

func taskFetchBankAccounts(logger logging.Logger, client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchBankAccounts)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), bankAccountsAttrs)
		}()

		for page := 1; ; page++ {
			pagedBankAccounts, err := client.GetBankAccounts(ctx, userID, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, bankAccountsAttrs)
				return err
			}

			if len(pagedBankAccounts) == 0 {
				break
			}

			var accountBatch ingestion.AccountBatch
			for _, bankAccount := range pagedBankAccounts {
				buf, err := json.Marshal(bankAccount)
				if err != nil {
					return err
				}

				accountBatch = append(accountBatch, &models.Account{
					ID: models.AccountID{
						Reference: bankAccount.ID,
						Provider:  models.ConnectorProviderMangopay,
					},
					CreatedAt:   time.Unix(bankAccount.CreationDate, 0),
					Reference:   bankAccount.ID,
					Provider:    models.ConnectorProviderMangopay,
					AccountName: bankAccount.OwnerName,
					Type:        models.AccountTypeExternal,
					Metadata: map[string]string{
						"user_id": userID,
					},
					RawData: buf,
				})
			}

			if err := ingester.IngestAccounts(ctx, accountBatch); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, bankAccountsAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountBatch)), bankAccountsAttrs)
		}

		return nil
	}
}
