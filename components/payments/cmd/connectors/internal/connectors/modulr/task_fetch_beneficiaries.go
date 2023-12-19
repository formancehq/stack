package modulr

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/internal/models"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchBeneficiaries(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		logger.Info(taskNameFetchBeneficiaries)

		beneficiaries, err := client.GetBeneficiaries(ctx)
		if err != nil {
			return err
		}

		if err := ingestBeneficiariesAccountsBatch(ctx, connectorID, ingester, beneficiaries); err != nil {
			return err
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
	accountsBatch := ingestion.AccountBatch{}

	for _, beneficiary := range beneficiaries {
		raw, err := json.Marshal(beneficiary)
		if err != nil {
			return err
		}

		openingDate, err := time.Parse("2006-01-02T15:04:05.999999999+0000", beneficiary.Created)
		if err != nil {
			return err
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference:   beneficiary.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:   openingDate,
			Reference:   beneficiary.ID,
			ConnectorID: connectorID,
			AccountName: beneficiary.Name,
			Type:        models.AccountTypeExternal,
			RawData:     raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		return err
	}

	return nil
}
