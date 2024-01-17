package mangopay

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func taskFetchBankAccounts(client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("mangopay.taskFetchBankAccounts")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("userID", userID),
		)

		if err := fetchBankAccounts(ctx, client, userID, connectorID, scheduler, ingester); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchBankAccounts(
	ctx context.Context,
	client *client.Client,
	userID string,
	connectorID models.ConnectorID,
	scheduler task.Scheduler,
	ingester ingestion.Ingester,
) error {

	for page := 1; ; page++ {
		pagedBankAccounts, err := client.GetBankAccounts(ctx, userID, page)
		if err != nil {
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
					Reference:   bankAccount.ID,
					ConnectorID: connectorID,
				},
				CreatedAt:   time.Unix(bankAccount.CreationDate, 0),
				Reference:   bankAccount.ID,
				ConnectorID: connectorID,
				AccountName: bankAccount.OwnerName,
				Type:        models.AccountTypeExternal,
				Metadata: map[string]string{
					"user_id": userID,
				},
				RawData: buf,
			})
		}

		if err := ingester.IngestAccounts(ctx, accountBatch); err != nil {
			return err
		}
	}

	return nil
}
