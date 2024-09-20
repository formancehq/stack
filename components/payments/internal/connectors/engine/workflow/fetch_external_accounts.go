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
	stateReference := models.CAPABILITY_FETCH_EXTERNAL_ACCOUNTS.String()
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

		accounts := models.FromPSPAccounts(
			externalAccountsResponse.ExternalAccounts,
			models.ACCOUNT_TYPE_EXTERNAL,
			fetchNextExternalAccount.ConnectorID,
		)

		if len(externalAccountsResponse.ExternalAccounts) > 0 {
			err = activities.StorageAccountsStore(
				infiniteRetryContext(ctx),
				accounts,
			)
			if err != nil {
				return errors.Wrap(err, "storing next accounts")
			}
		}

		wg := workflow.NewWaitGroup(ctx)
		errChan := make(chan error, len(externalAccountsResponse.ExternalAccounts)*2)
		for _, externalAccount := range accounts {
			acc := externalAccount
			wg.Add(1)
			workflow.Go(ctx, func(ctx workflow.Context) {
				defer wg.Done()

				if err := workflow.ExecuteChildWorkflow(
					workflow.WithChildOptions(
						ctx,
						workflow.ChildWorkflowOptions{
							TaskQueue:         fetchNextExternalAccount.ConnectorID.String(),
							ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
							SearchAttributes: map[string]interface{}{
								SearchAttributeStack: w.stack,
							},
						},
					),
					RunSendEvents,
					SendEvents{
						Account: &acc,
					},
				).Get(ctx, nil); err != nil {
					errChan <- errors.Wrap(err, "sending events")
				}
			})
		}

		for _, externalAccount := range externalAccountsResponse.ExternalAccounts {
			acc := externalAccount

			wg.Add(1)
			workflow.Go(ctx, func(ctx workflow.Context) {
				defer wg.Done()

				payload, err := json.Marshal(acc)
				if err != nil {
					errChan <- errors.Wrap(err, "marshalling external account")
				}

				if err := workflow.ExecuteChildWorkflow(
					workflow.WithChildOptions(
						ctx,
						workflow.ChildWorkflowOptions{
							TaskQueue:         fetchNextExternalAccount.ConnectorID.String(),
							ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
							SearchAttributes: map[string]interface{}{
								SearchAttributeStack: w.stack,
							},
						},
					),
					Run,
					fetchNextExternalAccount.Config,
					fetchNextExternalAccount.ConnectorID,
					&FromPayload{
						ID:      acc.Reference,
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

		state.State = externalAccountsResponse.NewState
		err = activities.StorageStatesStore(
			infiniteRetryContext(ctx),
			*state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		hasMore = externalAccountsResponse.HasMore
	}

	return nil
}

var RunFetchNextExternalAccounts any

func init() {
	RunFetchNextExternalAccounts = Workflow{}.runFetchNextExternalAccounts
}
