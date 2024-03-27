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
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

type fetchRecipientAccountsState struct {
	LastSeekPosition uint64 `json:"last_seek_position"`
}

func taskFetchRecipientAccounts(wiseClient *client.Client, profileID uint64) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"wise.taskFetchRecipientAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("profileID", strconv.FormatUint(profileID, 10)),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchRecipientAccountsState{})

		for {
			recipientAccounts, err := wiseClient.GetRecipientAccounts(ctx, profileID, pageSize, state.LastSeekPosition)
			if err != nil {
				// Retryable errors already handled by the function
				otel.RecordError(span, err)
				return err
			}

			if err := ingestRecipientAccountsBatch(ctx, connectorID, ingester, recipientAccounts.Content); err != nil {
				// Retryable errors already handled by the function
				otel.RecordError(span, err)
				return err
			}

			if recipientAccounts.SeekPositionForNext == 0 {
				// No more data to fetch
				break
			}

			state.LastSeekPosition = recipientAccounts.SeekPositionForNext
		}

		if err := ingester.UpdateTaskState(ctx, state); err != nil {
			otel.RecordError(span, err)
			return errors.Wrap(task.ErrRetryable, err.Error())
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
			return errors.Wrap(task.ErrRetryable, err.Error())
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
			AccountName:  account.Name.FullName,
			Type:         models.AccountTypeExternal,
			RawData:      raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		return errors.Wrap(task.ErrRetryable, err.Error())
	}

	return nil
}
