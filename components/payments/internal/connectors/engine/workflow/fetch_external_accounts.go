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

type FetchNextExternalAccounts struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload *FromPayload       `json:"fromPayload"`
}

func (w Workflow) runFetchNextExternalAccounts(
	ctx workflow.Context,
	fetchNextExternalAccount FetchNextExternalAccounts,
	nextTasks []models.TaskTree,
) error {
	if err := w.createInstance(ctx, fetchNextExternalAccount.ConnectorID); err != nil {
		return errors.Wrap(err, "creating instance")
	}
	err := w.fetchExternalAccounts(ctx, fetchNextExternalAccount, nextTasks)
	return w.terminateInstance(ctx, fetchNextExternalAccount.ConnectorID, err)
}

func (w Workflow) fetchExternalAccounts(
	ctx workflow.Context,
	fetchNextExternalAccount FetchNextExternalAccounts,
	nextTasks []models.TaskTree,
) error {
	stateReference := fmt.Sprintf("%s", models.CAPABILITY_FETCH_EXTERNAL_ACCOUNTS.String())
	if fetchNextExternalAccount.FromPayload != nil {
		stateReference = fmt.Sprintf("%s-%s", models.CAPABILITY_FETCH_EXTERNAL_ACCOUNTS.String(), fetchNextExternalAccount.FromPayload.ID)
	}

	stateID := models.StateID{
		Reference:   stateReference,
		ConnectorID: fetchNextExternalAccount.ConnectorID,
	}
	state, err := activities.StorageStatesGet(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return fmt.Errorf("retrieving state %s: %v", stateID.String(), err)
	}

	hasMore := true
	for hasMore {
		externalAccountsResponse, err := activities.PluginFetchNextExternalAccounts(
			infiniteRetryContext(ctx),
			fetchNextExternalAccount.ConnectorID,
			fetchNextExternalAccount.FromPayload.GetPayload(),
			state.State,
			fetchNextExternalAccount.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		if len(externalAccountsResponse.ExternalAccounts) > 0 {
			err = activities.StorageAccountsStore(
				infiniteRetryContext(ctx),
				models.FromPSPAccounts(
					externalAccountsResponse.ExternalAccounts,
					models.ACCOUNT_TYPE_EXTERNAL,
					fetchNextExternalAccount.ConnectorID,
				),
			)
			if err != nil {
				return errors.Wrap(err, "storing next accounts")
			}
		}

		state.State = externalAccountsResponse.NewState
		err = activities.StorageStatesStore(
			infiniteRetryContext(ctx),
			*state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		// TODO(polo): send event

		for _, externalAccount := range externalAccountsResponse.ExternalAccounts {
			payload, err := json.Marshal(externalAccount)
			if err != nil {
				return errors.Wrap(err, "marshalling external account")
			}

			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         fetchNextExternalAccount.ConnectorID.Reference,
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					},
				),
				Run,
				fetchNextExternalAccount.Config,
				fetchNextExternalAccount.ConnectorID,
				&FromPayload{
					ID:      externalAccount.Reference,
					Payload: payload,
				},
				nextTasks,
			).Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = externalAccountsResponse.HasMore
	}

	return nil
}

var RunFetchNextExternalAccounts any

func init() {
	RunFetchNextExternalAccounts = Workflow{}.runFetchNextExternalAccounts
}
