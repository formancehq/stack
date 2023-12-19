package mangopay

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchBankAccounts(logger logging.Logger, client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		logger.Info(taskNameFetchBankAccounts)

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
}
