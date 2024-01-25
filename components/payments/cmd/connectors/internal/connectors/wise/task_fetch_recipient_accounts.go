package wise

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

func taskFetchRecipientAccounts(wiseClient *client.Client, profileID uint64) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"wise.taskFetchRecipientAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("profileID", strconv.FormatUint(profileID, 10)),
		)
		defer span.End()

		recipientAccounts, err := wiseClient.GetRecipientAccounts(ctx, profileID)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingestRecipientAccountsBatch(ctx, connectorID, ingester, recipientAccounts); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func ingestRecipientAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
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
				Reference:   fmt.Sprintf("%d", account.ID),
				ConnectorID: connectorID,
			},
			CreatedAt:    time.Now(),
			Reference:    fmt.Sprintf("%d", account.ID),
			ConnectorID:  connectorID,
			DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency),
			AccountName:  account.HolderName,
			Type:         models.AccountTypeExternal,
			RawData:      raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		return err
	}

	return nil
}
