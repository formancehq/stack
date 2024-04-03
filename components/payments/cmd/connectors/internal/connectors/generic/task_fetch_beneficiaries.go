package generic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/generic/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

type fetchBeneficiariesState struct {
	LastCreatedAt time.Time
}

func taskFetchBeneficiaries(client *client.Client, config *Config) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"generic.taskFetchBeneficiaries",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchBeneficiariesState{})

		newState, err := ingestBeneficiaries(ctx, connectorID, client, ingester, state)
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

func ingestBeneficiaries(
	ctx context.Context,
	connectorID models.ConnectorID,
	client *client.Client,
	ingester ingestion.Ingester,
	state fetchBeneficiariesState,
) (fetchBeneficiariesState, error) {
	newState := fetchBeneficiariesState{
		LastCreatedAt: state.LastCreatedAt,
	}

	for page := 1; ; page++ {
		beneficiaries, err := client.ListBeneficiaries(ctx, int64(page), pageSize, state.LastCreatedAt)
		if err != nil {
			return fetchBeneficiariesState{}, err
		}

		if len(beneficiaries) == 0 {
			break
		}

		beneficiaryBatch := make([]*models.Account, 0, len(beneficiaries))
		for _, beneficiary := range beneficiaries {
			raw, err := json.Marshal(beneficiary)
			if err != nil {
				return fetchBeneficiariesState{}, err
			}

			beneficiaryBatch = append(beneficiaryBatch, &models.Account{
				ID: models.AccountID{
					Reference:   beneficiary.Id,
					ConnectorID: connectorID,
				},
				ConnectorID: connectorID,
				CreatedAt:   beneficiary.CreatedAt,
				Reference:   beneficiary.Id,
				AccountName: beneficiary.OwnerName,
				Type:        models.AccountTypeExternal,
				Metadata:    beneficiary.Metadata,
				RawData:     raw,
			})

			newState.LastCreatedAt = beneficiary.CreatedAt
		}

		if err := ingester.IngestAccounts(ctx, ingestion.AccountBatch(beneficiaryBatch)); err != nil {
			return fetchBeneficiariesState{}, errors.Wrap(task.ErrRetryable, err.Error())
		}
	}

	return newState, nil
}
