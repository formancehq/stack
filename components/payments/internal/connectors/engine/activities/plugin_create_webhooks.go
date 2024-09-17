package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

type CreateWebhooksRequest struct {
	ConnectorID models.ConnectorID
	Req         models.CreateWebhooksRequest
}

func (a Activities) PluginCreateWebhooks(ctx context.Context, request CreateWebhooksRequest) (*models.CreateWebhooksResponse, error) {
	plugin, err := a.plugins.Get(request.ConnectorID)
	if err != nil {
		return nil, err
	}

	resp, err := plugin.CreateWebhooks(ctx, request.Req)
	if err != nil {
		// TODO(polo): temporal errors
		return nil, err
	}
	return &resp, nil
}

var PluginCreateWebhooksActivity = Activities{}.PluginCreateWebhooks

func PluginCreateWebhooks(ctx workflow.Context, connectorID models.ConnectorID, request models.CreateWebhooksRequest) (*models.CreateWebhooksResponse, error) {
	ret := models.CreateWebhooksResponse{}
	if err := executeActivity(ctx, PluginCreateWebhooksActivity, &ret, CreateWebhooksRequest{
		ConnectorID: connectorID,
		Req:         request,
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}
