package mangopay

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchBankAccountsState struct {
	LastPage         int       `json:"last_page"`
	LastCreationDate time.Time `json:"last_creation_date"`
}

func taskFetchBankAccounts(client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"mangopay.taskFetchBankAccounts",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("userID", userID),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchBankAccountsState{})
		if state.LastPage == 0 {
			// If last page is 0, it means we haven't fetched any wallets yet.
			// Mangopay pages starts at 1.
			state.LastPage = 1
		}

		newState, err := ingestBankAccounts(ctx, client, userID, connectorID, scheduler, ingester, state)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingester.UpdateTaskState(ctx, newState); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func ingestBankAccounts(
	ctx context.Context,
	client *client.Client,
	userID string,
	connectorID models.ConnectorID,
	scheduler task.Scheduler,
	ingester ingestion.Ingester,
	state fetchBankAccountsState,
) (fetchBankAccountsState, error) {
	var currentPage int

	newState := fetchBankAccountsState{
		LastCreationDate: state.LastCreationDate,
	}

	for currentPage = state.LastPage; ; currentPage++ {
		pagedBankAccounts, err := client.GetBankAccounts(ctx, userID, currentPage, pageSize)
		if err != nil {
			return fetchBankAccountsState{}, err
		}

		if len(pagedBankAccounts) == 0 {
			break
		}

		var accountBatch ingestion.AccountBatch
		for _, bankAccount := range pagedBankAccounts {
			creationDate := time.Unix(bankAccount.CreationDate, 0)
			switch creationDate.Compare(state.LastCreationDate) {
			case -1, 0:
				// creationDate <= state.LastCreationDate, nothing to do,
				// we already processed this bank account.
				continue
			default:
			}
			newState.LastCreationDate = creationDate

			buf, err := json.Marshal(bankAccount)
			if err != nil {
				return fetchBankAccountsState{}, err
			}

			accountBatch = append(accountBatch, &models.Account{
				ID: models.AccountID{
					Reference:   bankAccount.ID,
					ConnectorID: connectorID,
				},
				CreatedAt:   creationDate,
				Reference:   bankAccount.ID,
				ConnectorID: connectorID,
				AccountName: bankAccount.OwnerName,
				Type:        models.AccountTypeExternal,
				Metadata: map[string]string{
					"user_id": userID,
				},
				RawData: buf,
			})

			newState.LastCreationDate = creationDate
		}

		if err := ingester.IngestAccounts(ctx, accountBatch); err != nil {
			return fetchBankAccountsState{}, err
		}

		if len(pagedBankAccounts) < pageSize {
			break
		}
	}

	newState.LastPage = currentPage

	return newState, nil
}
