package workflow

import (
	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type CreateWebhooks struct {
	ConnectorID models.ConnectorID
	Config      models.Config
	FromPayload *FromPayload
}

func (w Workflow) runCreateWebhooks(
	ctx workflow.Context,
	createWebhooks CreateWebhooks,
	nextTasks []models.TaskTree,
) error {
	if err := w.createInstance(ctx, createWebhooks.ConnectorID); err != nil {
		return errors.Wrap(err, "creating instance")
	}
	err := w.createWebhooks(ctx, createWebhooks, nextTasks)
	return w.terminateInstance(ctx, createWebhooks.ConnectorID, err)
}

func (w Workflow) createWebhooks(
	ctx workflow.Context,
	createWebhooks CreateWebhooks,
	nextTasks []models.TaskTree,
) error {
	resp, err := activities.PluginCreateWebhooks(
		infiniteRetryContext(ctx),
		createWebhooks.ConnectorID,
		models.CreateWebhooksRequest{
			ConnectorID: createWebhooks.ConnectorID.String(),
			FromPayload: createWebhooks.FromPayload.GetPayload(),
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to create webhooks")
	}

	for _, other := range resp.Others {
		if err := workflow.ExecuteChildWorkflow(
			workflow.WithChildOptions(
				ctx,
				workflow.ChildWorkflowOptions{
					TaskQueue:         createWebhooks.ConnectorID.String(),
					ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					SearchAttributes: map[string]interface{}{
						SearchAttributeStack: w.stack,
					},
				},
			),
			Run,
			createWebhooks.Config,
			createWebhooks.ConnectorID,
			&FromPayload{
				ID:      other.ID,
				Payload: other.Other,
			},
			nextTasks,
		).Get(ctx, nil); err != nil {
			return errors.Wrap(err, "running next workflow")
		}
	}

	return nil
}

var RunCreateWebhooks any

func init() {
	RunCreateWebhooks = Workflow{}.runCreateWebhooks
}
