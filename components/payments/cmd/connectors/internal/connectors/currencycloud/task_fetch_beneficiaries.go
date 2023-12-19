package currencycloud

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchBeneficiaries(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		logger.Info(taskFetchBeneficiaries)

		page := 1
		for {
			if page < 0 {
				break
			}

			pagedBeneficiaries, nextPage, err := client.GetBeneficiaries(ctx, page)
			if err != nil {
				return err
			}

			page = nextPage

			if err := ingestBeneficiariesAccountsBatch(ctx, connectorID, ingester, pagedBeneficiaries); err != nil {
				return err
			}
		}

		return nil
	}
}

func ingestBeneficiariesAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	beneficiaries []*client.Beneficiary,
) error {
	batch := ingestion.AccountBatch{}
	for _, beneficiary := range beneficiaries {
		raw, err := json.Marshal(beneficiary)
		if err != nil {
			return err
		}

		batch = append(batch, &models.Account{
			ID: models.AccountID{
				Reference:   beneficiary.ID,
				ConnectorID: connectorID,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:   beneficiary.CreatedAt,
			Reference:   beneficiary.ID,
			ConnectorID: connectorID,
			AccountName: beneficiary.Name,
			Type:        models.AccountTypeExternal,
			RawData:     raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
