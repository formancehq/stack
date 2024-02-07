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

type fetchBeneficiariesState struct {
	LastCreated time.Time `json:"last_created"`
}

func (s *fetchBeneficiariesState) SetLatestBeneficiary(latest *client.Beneficiary) error {
	createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", latest.Created)
	if err != nil {
		return err
	}
	if createdTime.After(s.LastCreated) {
		s.LastCreated = createdTime
	}
	return nil
}

func (s *fetchBeneficiariesState) GetLastCreated() string {
	if s.LastCreated.IsZero() {
		return ""
	}
	return s.LastCreated.Format("2006-01-02T15:04:05-0700")
}

func taskFetchBeneficiaries(config Config, client *client.Client) task.Task {
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
			"modulr.taskFetchBeneficiaries",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state, err := fetchBeneficiaries(
			ctx,
			config,
			client,
			connectorID,
			ingester,
			scheduler,
			task.MustResolveTo(ctx, resolver, fetchBeneficiariesState{}),
		)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingester.UpdateTaskState(ctx, state); err != nil {
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
	state fetchBeneficiariesState,
) (fetchBeneficiariesState, error) {
	fromCreateTime := state.GetLastCreated()
	for page := 0; ; page++ {
		pagedBeneficiaries, err := client.GetBeneficiaries(ctx, page, config.PageSize, fromCreateTime)
		if err != nil {
			return state, err
		}

		if len(pagedBeneficiaries.Content) == 0 {
			break
		}

		for _, beneficiary := range pagedBeneficiaries.Content {
			if err := state.SetLatestBeneficiary(beneficiary); err != nil {
				return state, err
			}
		}

		if err := ingestBeneficiariesAccountsBatch(ctx, connectorID, ingester, pagedBeneficiaries.Content); err != nil {
			return state, err
		}

		if len(pagedBeneficiaries.Content) < config.PageSize {
			break
		}

		if page+1 >= pagedBeneficiaries.TotalPages {
			// Modulr paging starts at 0, so the last page is TotalPages - 1.
			break
		}
	}

	return state, nil
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
