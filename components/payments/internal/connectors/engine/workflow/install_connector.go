package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type InstallConnector struct {
	ConnectorID models.ConnectorID
	RawConfig   json.RawMessage
}

func (w Workflow) runInstallConnector(
	ctx workflow.Context,
	installConnector InstallConnector,
) error {
	// First step: store the connector inside the database
	connector := models.Connector{
		ID:        installConnector.ConnectorID,
		Name:      installConnector.ConnectorID.Reference,
		CreatedAt: workflow.Now(ctx).UTC(),
		Provider:  installConnector.ConnectorID.Provider,
		Config:    installConnector.RawConfig,
	}
	err := activities.StorageConnectorsStore(infiniteRetryContext(ctx), connector)
	if err != nil {
		return errors.Wrap(err, "failed to store connector")
	}

	// Second step: install the connector via the plugin and get the list of
	// capabilities and the workflow of polling data
	installResponse, err := activities.PluginInstallConnector(
		infiniteRetryContext(ctx),
		installConnector.ConnectorID,
		installConnector.RawConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to install connector")
	}

	// Third step: store the workflow of the connector
	err = activities.StorageTasksTreeStore(infiniteRetryContext(ctx), installConnector.ConnectorID, installResponse.Workflow)
	if err != nil {
		return errors.Wrap(err, "failed to store tasks tree")
	}

	configs := make([]models.WebhookConfig, 0, len(installResponse.WebhooksConfigs))
	for _, webhookConfig := range installResponse.WebhooksConfigs {
		configs = append(configs, models.WebhookConfig{
			Name:        webhookConfig.Name,
			ConnectorID: installConnector.ConnectorID,
			URLPath:     webhookConfig.URLPath,
		})
	}
	// TODO(polo): store the capabilities
	err = activities.StorageWebhooksConfigsStore(infiniteRetryContext(ctx), configs)
	if err != nil {
		return errors.Wrap(err, "failed to store webhooks configs")
	}

	var config models.Config
	if err := json.Unmarshal(installConnector.RawConfig, &config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	// Fourth step: launch the workflow tree
	if err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(
			ctx,
			workflow.ChildWorkflowOptions{
				WorkflowID:            fmt.Sprintf("run-tasks-%s", installConnector.ConnectorID.String()),
				WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
				TaskQueue:             installConnector.ConnectorID.Reference,
				ParentClosePolicy:     enums.PARENT_CLOSE_POLICY_ABANDON,
			},
		),
		Run,
		config,
		installConnector.ConnectorID,
		nil,
		[]models.TaskTree(installResponse.Workflow),
	).Get(ctx, nil); err != nil {
		applicationError := &temporal.ApplicationError{}
		if errors.As(err, &applicationError) {
			if applicationError.Type() != "ChildWorkflowExecutionAlreadyStartedError" {
				return err
			}
		} else {
			return errors.Wrap(err, "running next workflow")
		}
	}

	return nil
}

var RunInstallConnector any

func init() {
	RunInstallConnector = Workflow{}.runInstallConnector
}