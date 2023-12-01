package bankingcircle

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	bankAccountCreationAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "bank_account"))...)
)

// No need to call any API for banking circle since it does not support it.
// We will just create an external accounts on our side linked to the
// bank account object.
func taskCreateExternalAccount(
	logger logging.Logger,
	client *client.Client,
	bankAccountID uuid.UUID,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		storageReader storage.Reader,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		bankAccount, err := storageReader.GetBankAccount(ctx, bankAccountID, false)
		if err != nil {
			return err
		}

		accountID := models.AccountID{
			Reference:   bankAccount.ID.String(),
			ConnectorID: connectorID,
		}

		if err := ingester.IngestAccounts(ctx, ingestion.AccountBatch{
			{
				ID:          accountID,
				CreatedAt:   time.Now(),
				Reference:   bankAccount.ID.String(),
				ConnectorID: connectorID,
				AccountName: bankAccount.Name,
				Type:        models.AccountTypeExternalFormance,
			},
		}); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, bankAccountCreationAttrs)
			return err
		}

		if err = ingester.LinkBankAccountWithAccount(ctx, bankAccount, &accountID); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, bankAccountCreationAttrs)
			return err
		}

		metricsRegistry.ConnectorObjects().Add(ctx, 1, bankAccountCreationAttrs)
		return nil
	}
}
