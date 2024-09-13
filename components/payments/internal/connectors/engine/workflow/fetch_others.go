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

type FetchNextOthers struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	Name        string             `json:"name"`
	FromPayload *FromPayload       `json:"fromPayload"`
}

func (w Workflow) runFetchNextOthers(
	ctx workflow.Context,
	fetchNextOthers FetchNextOthers,
	nextTasks []models.TaskTree,
) error {
	if err := w.createInstance(ctx, fetchNextOthers.ConnectorID); err != nil {
		return errors.Wrap(err, "creating instance")
	}
	err := w.fetchNextOthers(ctx, fetchNextOthers, nextTasks)
	return w.terminateInstance(ctx, fetchNextOthers.ConnectorID, err)
}

func (w Workflow) fetchNextOthers(
	ctx workflow.Context,
	fetchNextOthers FetchNextOthers,
	nextTasks []models.TaskTree,
) error {
	stateReference := fmt.Sprintf("%s", models.CAPABILITY_FETCH_OTHERS.String())
	if fetchNextOthers.FromPayload != nil {
		stateReference = fmt.Sprintf("%s-%s", models.CAPABILITY_FETCH_OTHERS.String(), fetchNextOthers.FromPayload.ID)
	}

	stateID := models.StateID{
		Reference:   stateReference,
		ConnectorID: fetchNextOthers.ConnectorID,
	}
	state, err := activities.StorageStatesGet(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return fmt.Errorf("retrieving state %s: %v", stateID.String(), err)
	}

	hasMore := true
	for hasMore {
		othersResponse, err := activities.PluginFetchNextOthers(
			infiniteRetryContext(ctx),
			fetchNextOthers.ConnectorID,
			fetchNextOthers.Name,
			fetchNextOthers.FromPayload.GetPayload(),
			state.State,
			fetchNextOthers.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next others")
		}

		state.State = othersResponse.NewState
		err = activities.StorageStatesStore(
			infiniteRetryContext(ctx),
			*state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		// TODO(polo): send event for others ? store others ?

		for _, other := range othersResponse.Others {
			payload, err := json.Marshal(other.Other)
			if err != nil {
				return errors.Wrap(err, "marshalling other")
			}

			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         fetchNextOthers.ConnectorID.Reference,
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					},
				),
				Run,
				fetchNextOthers.Config,
				fetchNextOthers.ConnectorID,
				&FromPayload{
					ID:      other.ID,
					Payload: payload,
				},
				nextTasks,
			).Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = othersResponse.HasMore
	}

	return nil
}

var RunFetchNextOthers any

func init() {
	RunFetchNextOthers = Workflow{}.runFetchNextOthers
}
