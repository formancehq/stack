package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

type FetchNextOthersRequest struct {
	ConnectorID models.ConnectorID
	Req         models.FetchNextOthersRequest
}

func (a Activities) PluginFetchNextOthers(ctx context.Context, request FetchNextOthersRequest) (*models.FetchNextOthersResponse, error) {
	plugin, err := a.plugins.Get(request.ConnectorID)
	if err != nil {
		return nil, err
	}

	resp, err := plugin.FetchNextOthers(ctx, request.Req)
	if err != nil {
		// TODO(polo): temporal errors
		return nil, err
	}

	return &resp, nil
}

var PluginFetchNextOthersActivity = Activities{}.PluginFetchNextOthers

func PluginFetchNextOthers(ctx workflow.Context, connectorID models.ConnectorID, name string, fromPayload, state json.RawMessage, pageSize int) (*models.FetchNextOthersResponse, error) {
	ret := models.FetchNextOthersResponse{}
	if err := executeActivity(ctx, PluginFetchNextOthersActivity, &ret, FetchNextOthersRequest{
		ConnectorID: connectorID,
		Req: models.FetchNextOthersRequest{
			Name:        name,
			FromPayload: fromPayload,
			State:       state,
			PageSize:    pageSize,
		},
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}
