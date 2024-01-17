package bankingcircle

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// No need to call any API for banking circle since it does not support it.
// We will just create an external accounts on our side linked to the
// bank account object.
func taskCreateExternalAccount(
	client *client.Client,
	bankAccountID uuid.UUID,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		storageReader storage.Reader,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("bankingcircle.taskCreateExternalAccount")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("bankAccount.id", bankAccountID.String()),
		)

		bankAccount, err := storageReader.GetBankAccount(ctx, bankAccountID, false)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		span.SetAttributes(attribute.String("bankAccount.name", bankAccount.Name))

		if err := createExternalAccount(ctx, connectorID, ingester, storageReader, bankAccount); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func createExternalAccount(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	storageReader storage.Reader,
	bankAccount *models.BankAccount,
) error {
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
		return err
	}

	if err := ingester.LinkBankAccountWithAccount(ctx, bankAccount, &accountID); err != nil {
		return err
	}

	return nil
}
