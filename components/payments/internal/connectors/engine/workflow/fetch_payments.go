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
	stateReference := models.CAPABILITY_FETCH_PAYMENTS.String()
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

		payments := models.FromPSPPayments(
			paymentsResponse.Payments,
			fetchNextPayments.ConnectorID,
		)

		if len(paymentsResponse.Payments) > 0 {
			err = activities.StoragePaymentsStore(
				infiniteRetryContext(ctx),
				payments,
			)
			if err != nil {
				return errors.Wrap(err, "storing next accounts")
			}
		}

		wg := workflow.NewWaitGroup(ctx)
		errChan := make(chan error, len(paymentsResponse.Payments)*2)
		for _, payment := range payments {
			p := payment

			wg.Add(1)
			workflow.Go(ctx, func(ctx workflow.Context) {
				defer wg.Done()

				if err := workflow.ExecuteChildWorkflow(
					workflow.WithChildOptions(
						ctx,
						workflow.ChildWorkflowOptions{
							TaskQueue:         fetchNextPayments.ConnectorID.String(),
							ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
							SearchAttributes: map[string]interface{}{
								SearchAttributeStack: w.stack,
							},
						},
					),
					RunSendEvents,
					SendEvents{
						Payment: &p,
					},
				).Get(ctx, nil); err != nil {
					errChan <- errors.Wrap(err, "sending events")
				}
			})
		}

		for _, payment := range paymentsResponse.Payments {
			p := payment

			wg.Add(1)
			workflow.Go(ctx, func(ctx workflow.Context) {
				defer wg.Done()

				payload, err := json.Marshal(p)
				if err != nil {
					errChan <- errors.Wrap(err, "marshalling payment")
				}

				if err := workflow.ExecuteChildWorkflow(
					workflow.WithChildOptions(
						ctx,
						workflow.ChildWorkflowOptions{
							TaskQueue:         fetchNextPayments.ConnectorID.String(),
							ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
							SearchAttributes: map[string]interface{}{
								SearchAttributeStack: w.stack,
							},
						},
					),
					Run,
					fetchNextPayments.Config,
					fetchNextPayments.ConnectorID,
					&FromPayload{
						ID:      p.Reference,
						Payload: payload,
					},
					nextTasks,
				).Get(ctx, nil); err != nil {
					errChan <- errors.Wrap(err, "running next workflow")
				}
			})
		}

		wg.Wait(ctx)
		close(errChan)
		for err := range errChan {
			if err != nil {
				return err
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

		hasMore = paymentsResponse.HasMore
	}

	return nil
}

var RunFetchNextPayments any

func init() {
	RunFetchNextPayments = Workflow{}.runFetchNextPayments
}
