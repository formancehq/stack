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

func (s *fetchBeneficiariesState) UpdateLatest(latest *client.Beneficiary) error {
	createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", latest.Created)
	if err != nil {
		return err
	}
	if createdTime.After(s.LastCreated) {
		s.LastCreated = createdTime
	}
	return nil
}

func (s *fetchBeneficiariesState) FindLatest(beneficiaries []*client.Beneficiary) error {
	for _, beneficiary := range beneficiaries {
		if err := s.UpdateLatest(beneficiary); err != nil {
			return err
		}
	}
	return nil
}

func (s *fetchBeneficiariesState) IsNew(beneficiary *client.Beneficiary) (bool, error) {
	createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", beneficiary.Created)
	if err != nil {
		return false, err
	}
	return createdTime.After(s.LastCreated), nil
}

func (s *fetchBeneficiariesState) FilterNew(beneficiaries []*client.Beneficiary) ([]*client.Beneficiary, error) {
	// beneficiaries are not assumed to be sorted by creation date.
	result := make([]*client.Beneficiary, 0, len(beneficiaries))
	for _, beneficiary := range beneficiaries {
		isNew, err := s.IsNew(beneficiary)
		if err != nil {
			return nil, err
		}
		if !isNew {
			continue
		}
		result = append(result, beneficiary)
	}
	return result, nil
}

func (s *fetchBeneficiariesState) GetFilterValue() string {
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
	newState := state
	for page := 0; ; page++ {
		pagedBeneficiaries, err := client.GetBeneficiaries(
			ctx,
			page,
			config.PageSize,
			state.GetFilterValue(),
		)
		if err != nil {
			return newState, err
		}
		if len(pagedBeneficiaries.Content) == 0 {
			break
		}
		beneficiaries, err := state.FilterNew(pagedBeneficiaries.Content)
		if err != nil {
			return newState, err
		}
		if err := newState.FindLatest(beneficiaries); err != nil {
			return newState, err
		}
		if err := ingestBeneficiariesAccountsBatch(ctx, connectorID, ingester, beneficiaries); err != nil {
			return newState, err
		}

		if len(pagedBeneficiaries.Content) < config.PageSize {
			break
		}

		if page+1 >= pagedBeneficiaries.TotalPages {
			// Modulr paging starts at 0, so the last page is TotalPages - 1.
			break
		}
	}

	return newState, nil
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
