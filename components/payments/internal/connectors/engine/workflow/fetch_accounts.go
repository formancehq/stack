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

type FetchNextAccounts struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload *FromPayload       `json:"fromPayload"`
}

func (w Workflow) runFetchNextAccounts(
	ctx workflow.Context,
	fetchNextAccount FetchNextAccounts,
	nextTasks []models.TaskTree,
) error {
	if err := w.createInstance(ctx, fetchNextAccount.ConnectorID); err != nil {
		return errors.Wrap(err, "creating instance")
	}
	err := w.fetchAccounts(ctx, fetchNextAccount, nextTasks)
	return w.terminateInstance(ctx, fetchNextAccount.ConnectorID, err)
}

func (w Workflow) fetchAccounts(
	ctx workflow.Context,
	fetchNextAccount FetchNextAccounts,
	nextTasks []models.TaskTree,
) error {
	stateReference := fmt.Sprintf("%s", models.CAPABILITY_FETCH_ACCOUNTS.String())
	if fetchNextAccount.FromPayload != nil {
		stateReference = fmt.Sprintf("%s-%s", models.CAPABILITY_FETCH_ACCOUNTS.String(), fetchNextAccount.FromPayload.ID)
	}

	stateID := models.StateID{
		Reference:   stateReference,
		ConnectorID: fetchNextAccount.ConnectorID,
	}
	state, err := activities.StorageStatesGet(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return fmt.Errorf("retrieving state %s: %v", stateID.String(), err)
	}

	hasMore := true
	for hasMore {
		accountsResponse, err := activities.PluginFetchNextAccounts(
			infiniteRetryContext(ctx),
			fetchNextAccount.ConnectorID,
			fetchNextAccount.FromPayload.GetPayload(),
			state.State,
			fetchNextAccount.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		if len(accountsResponse.Accounts) > 0 {
			err = activities.StorageAccountsStore(
				infiniteRetryContext(ctx),
				models.FromPSPAccounts(
					accountsResponse.Accounts,
					models.ACCOUNT_TYPE_INTERNAL,
					fetchNextAccount.ConnectorID,
				),
			)
			if err != nil {
				return errors.Wrap(err, "storing next accounts")
			}
		}

		// TODO(polo): send event
		// TODO(polo): events with IK to avoid duplicates
		for _, account := range accountsResponse.Accounts {
			payload, err := json.Marshal(account)
			if err != nil {
				return errors.Wrap(err, "marshalling account")
			}

			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         fetchNextAccount.ConnectorID.String(),
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					},
				),
				Run,
				fetchNextAccount.Config,
				fetchNextAccount.ConnectorID,
				&FromPayload{
					ID:      account.Reference,
					Payload: payload,
				},
				nextTasks,
			).Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		state.State = accountsResponse.NewState
		err = activities.StorageStatesStore(
			infiniteRetryContext(ctx),
			*state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		hasMore = accountsResponse.HasMore
	}

	return nil
}

var RunFetchNextAccounts any

func init() {
	RunFetchNextAccounts = Workflow{}.runFetchNextAccounts
}
