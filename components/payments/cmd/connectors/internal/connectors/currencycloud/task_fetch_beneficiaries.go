package currencycloud

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchBeneficiariesState struct {
	LastPage      int
	LastCreatedAt time.Time
}

func taskFetchBeneficiaries(
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"currencycloud.taskFetchBeneficiaries",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchBeneficiariesState{})
		if state.LastPage == 0 {
			// First run, the first page for currencycloud starts at 1 and not 0
			state.LastPage = 1
		}

		newState, err := fetchBeneficiaries(ctx, client, connectorID, ingester, scheduler, state)
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

func fetchBeneficiaries(
	ctx context.Context,
	client *client.Client,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	state fetchBeneficiariesState,
) (fetchBeneficiariesState, error) {
	newState := fetchBeneficiariesState{
		LastPage:      state.LastPage,
		LastCreatedAt: state.LastCreatedAt,
	}

	page := state.LastPage
	for {
		if page < 0 {
			break
		}

		pagedBeneficiaries, nextPage, err := client.GetBeneficiaries(ctx, page)
		if err != nil {
			return fetchBeneficiariesState{}, err
		}

		page = nextPage

		batch := ingestion.AccountBatch{}
		for _, beneficiary := range pagedBeneficiaries {
			switch beneficiary.CreatedAt.Compare(state.LastCreatedAt) {
			case -1, 0:
				// Account already ingested, skip
				continue
			default:
			}

			raw, err := json.Marshal(beneficiary)
			if err != nil {
				return fetchBeneficiariesState{}, err
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

			newState.LastCreatedAt = beneficiary.CreatedAt
		}

		if err := ingester.IngestAccounts(ctx, batch); err != nil {
			return fetchBeneficiariesState{}, err
		}
	}

	newState.LastPage = page

	return newState, nil
}
