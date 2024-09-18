package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

type UninstallConnectorRequest struct {
	ConnectorID models.ConnectorID
}

func (a Activities) PluginUninstallConnector(ctx context.Context, request UninstallConnectorRequest) (*models.UninstallResponse, error) {
	plugin, err := a.plugins.Get(request.ConnectorID)
	if err != nil {
		return nil, err
	}

	resp, err := plugin.Uninstall(ctx, models.UninstallRequest{})
	if err != nil {
		// TODO(polo): temporal errors
		return nil, err
	}

	return &resp, err
}

var PluginUninstallConnectorActivity = Activities{}.PluginUninstallConnector

func PluginUninstallConnector(ctx workflow.Context, connectorID models.ConnectorID) (*models.UninstallResponse, error) {
	ret := models.UninstallResponse{}
	if err := executeActivity(ctx, PluginUninstallConnectorActivity, &ret, UninstallConnectorRequest{
		ConnectorID: connectorID,
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}
