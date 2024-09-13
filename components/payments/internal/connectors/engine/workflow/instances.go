package workflow

import (
	"encoding/json"

	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

func (w Workflow) createInstance(
	ctx workflow.Context,
	connectorID models.ConnectorID,
) error {
	info := workflow.GetInfo(ctx)

	scheduleID, err := getPaymentScheduleID(ctx, info)
	if err != nil {
		return err
	}

	instance := models.Instance{
		ID:          info.WorkflowExecution.ID,
		ScheduleID:  scheduleID,
		ConnectorID: connectorID,
		CreatedAt:   workflow.Now(ctx).UTC(),
		UpdatedAt:   workflow.Now(ctx).UTC(),
		Terminated:  false,
	}

	return activities.StorageInstancesStore(infiniteRetryContext(ctx), instance)
}

func (w Workflow) terminateInstance(
	ctx workflow.Context,
	connectorID models.ConnectorID,
	err error,
) error {
	info := workflow.GetInfo(ctx)

	scheduleID, err := getPaymentScheduleID(ctx, info)
	if err != nil {
		return err
	}

	var errMessage *string
	if err != nil {
		errMessage = pointer.For(err.Error())
	}

	now := workflow.Now(ctx).UTC()

	instance := models.Instance{
		ID:           info.WorkflowExecution.ID,
		ScheduleID:   scheduleID,
		ConnectorID:  connectorID,
		UpdatedAt:    now,
		Terminated:   true,
		TerminatedAt: &now,
		Error:        errMessage,
	}

	return activities.StorageInstancesUpdate(infiniteRetryContext(ctx), instance)
}

func getPaymentScheduleID(
	ctx workflow.Context,
	info *workflow.Info,
) (string, error) {
	attributes := info.SearchAttributes.GetIndexedFields()
	if attributes == nil {
		return "", errors.New("missing search attributes")
	}

	v, ok := attributes[SearchAttributeScheduleID]
	if !ok || v == nil {
		return "", errors.New("missing schedule ID")
	}

	var scheduleID string
	if err := json.Unmarshal(v.Data, &scheduleID); err != nil {
		return "", errors.Wrap(err, "unmarshalling schedule ID")
	}

	return scheduleID, nil
}
