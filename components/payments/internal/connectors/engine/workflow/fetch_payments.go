package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type FetchNextPayments struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload *FromPayload       `json:"fromPayload"`
}

func (w Workflow) runFetchNextPayments(
	ctx workflow.Context,
	fetchNextPayments FetchNextPayments,
	nextTasks []models.TaskTree,
) error {
	if err := w.createInstance(ctx, fetchNextPayments.ConnectorID); err != nil {
		return errors.Wrap(err, "creating instance")
	}
	err := w.fetchNextPayments(ctx, fetchNextPayments, nextTasks)
	return w.terminateInstance(ctx, fetchNextPayments.ConnectorID, err)
}

func (w Workflow) fetchNextPayments(
	ctx workflow.Context,
	fetchNextPayments FetchNextPayments,
	nextTasks []models.TaskTree,
) error {
	stateReference := fmt.Sprintf("%s", models.CAPABILITY_FETCH_PAYMENTS.String())
	if fetchNextPayments.FromPayload != nil {
		stateReference = fmt.Sprintf("%s-%s", models.CAPABILITY_FETCH_PAYMENTS.String(), fetchNextPayments.FromPayload.ID)
	}

	stateID := models.StateID{
		Reference:   stateReference,
		ConnectorID: fetchNextPayments.ConnectorID,
	}
	state, err := activities.StorageStatesGet(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return fmt.Errorf("retrieving state %s: %v", stateID.String(), err)
	}

	hasMore := true
	for hasMore {
		paymentsResponse, err := activities.PluginFetchNextPayments(
			infiniteRetryContext(ctx),
			fetchNextPayments.ConnectorID,
			fetchNextPayments.FromPayload.GetPayload(),
			state.State,
			fetchNextPayments.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next payments")
		}

		if len(paymentsResponse.Payments) > 0 {
			err = activities.StoragePaymentsStore(
				infiniteRetryContext(ctx),
				models.FromPSPPayments(
					paymentsResponse.Payments,
					fetchNextPayments.ConnectorID,
				),
			)
			if err != nil {
				return errors.Wrap(err, "storing next accounts")
			}
		}

		state.State = paymentsResponse.NewState
		err = activities.StorageStatesStore(
			infiniteRetryContext(ctx),
			*state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		// TODO(polo): send events

		for _, payment := range paymentsResponse.Payments {
			payload, err := json.Marshal(payment)
			if err != nil {
				return errors.Wrap(err, "marshalling payment")
			}

			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         fetchNextPayments.ConnectorID.Reference,
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					},
				),
				Run,
				fetchNextPayments.Config,
				fetchNextPayments.ConnectorID,
				&FromPayload{
					ID:      payment.Reference,
					Payload: payload,
				},
				nextTasks,
			).Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = paymentsResponse.HasMore
	}

	return nil
}

var RunFetchNextPayments any

func init() {
	RunFetchNextPayments = Workflow{}.runFetchNextPayments
}
