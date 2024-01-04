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

func taskFetchBeneficiaries(logger logging.Logger, config Config, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		logger.Info(taskNameFetchBeneficiaries)

		for page := 0; ; page++ {
			pagedBeneficiaries, err := client.GetBeneficiaries(ctx, page, config.PageSize)
			if err != nil {
				return err
			}

			if len(pagedBeneficiaries.Content) == 0 {
				break
			}

			if err := ingestBeneficiariesAccountsBatch(ctx, connectorID, ingester, pagedBeneficiaries.Content); err != nil {
				return err
			}

			if len(pagedBeneficiaries.Content) < config.PageSize {
				break
			}

			if page+1 >= pagedBeneficiaries.TotalPages {
				// Modulr paging starts at 0, so the last page is TotalPages - 1.
				break
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
