package moneycorp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

func taskFetchRecipients(client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"moneycorp.taskFetchRecipients",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("accountID", accountID),
		)
		defer span.End()

		if err := fetchRecipients(ctx, client, accountID, connectorID, ingester, scheduler); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchRecipients(
	ctx context.Context,
	client *client.Client,
	accountID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
) error {
	for page := 1; ; page++ {
		pagedRecipients, err := client.GetRecipients(ctx, accountID, page, pageSize)
		if err != nil {
			return err
		}

		if len(pagedRecipients) == 0 {
			break
		}

		if err := ingestRecipientsBatch(ctx, connectorID, ingester, pagedRecipients); err != nil {
			return err
		}

		if len(pagedRecipients) < pageSize {
			break
		}
	}

	return nil
}

func ingestRecipientsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
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
				Reference:   recipient.ID,
				ConnectorID: connectorID,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:    createdAt,
			Reference:    recipient.ID,
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, recipient.Attributes.BankAccountCurrency),
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
