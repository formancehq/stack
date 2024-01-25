package modulr

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
)

func taskFetchBeneficiaries(config Config, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"modulr.taskFetchBeneficiaries",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		if err := fetchBeneficiaries(ctx, config, client, connectorID, ingester, scheduler); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchBeneficiaries(
	ctx context.Context,
	config Config,
	client *client.Client,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
) error {
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

		openingDate, err := time.Parse(timeTemplate, beneficiary.Created)
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
