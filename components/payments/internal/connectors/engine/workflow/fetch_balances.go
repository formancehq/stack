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

type FetchNextBalances struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload *FromPayload       `json:"fromPayload"`
}

func (w Workflow) runFetchNextBalances(
	ctx workflow.Context,
	fetchNextBalances FetchNextBalances,
	nextTasks []models.TaskTree,
) error {
	if err := w.createInstance(ctx, fetchNextBalances.ConnectorID); err != nil {
		return errors.Wrap(err, "creating instance")
	}
	err := w.fetchBalances(ctx, fetchNextBalances, nextTasks)
	return w.terminateInstance(ctx, fetchNextBalances.ConnectorID, err)
}

func (w Workflow) fetchBalances(
	ctx workflow.Context,
	fetchNextBalances FetchNextBalances,
	nextTasks []models.TaskTree,
) error {
	stateReference := fmt.Sprintf("%s", models.CAPABILITY_FETCH_ACCOUNTS.String())
	if fetchNextBalances.FromPayload != nil {
		stateReference = fmt.Sprintf("%s-%s", models.CAPABILITY_FETCH_ACCOUNTS.String(), fetchNextBalances.FromPayload.ID)
	}

	stateID := models.StateID{
		Reference:   stateReference,
		ConnectorID: fetchNextBalances.ConnectorID,
	}
	state, err := activities.StorageStatesGet(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return fmt.Errorf("retrieving state %s: %v", stateID.String(), err)
	}

	hasMore := true
	for hasMore {
		balancesResponse, err := activities.PluginFetchNextBalances(
			infiniteRetryContext(ctx),
			fetchNextBalances.ConnectorID,
			fetchNextBalances.FromPayload.GetPayload(),
			state.State,
			fetchNextBalances.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		if len(balancesResponse.Balances) > 0 {
			err = activities.StorageBalancesStore(
				infiniteRetryContext(ctx),
				models.FromPSPBalances(
					balancesResponse.Balances,
					fetchNextBalances.ConnectorID,
				),
			)
			if err != nil {
				return errors.Wrap(err, "storing next accounts")
			}
		}

		state.State = balancesResponse.NewState
		err = activities.StorageStatesStore(
			infiniteRetryContext(ctx),
			*state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		// TODO(polo): send event

		for _, balance := range balancesResponse.Balances {
			payload, err := json.Marshal(balance)
			if err != nil {
				return errors.Wrap(err, "marshalling account")
			}

			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         fetchNextBalances.ConnectorID.Reference,
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					},
				),
				Run,
				fetchNextBalances.Config,
				fetchNextBalances.ConnectorID,
				&FromPayload{
					ID:      fmt.Sprintf("%s-balances", balance.AccountReference),
					Payload: payload,
				},
				nextTasks,
			).Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = balancesResponse.HasMore
	}

	return nil
}

var RunFetchNextBalances any

func init() {
	RunFetchNextBalances = Workflow{}.runFetchNextBalances
}