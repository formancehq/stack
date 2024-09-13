package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

type FetchNextBalancesRequest struct {
	ConnectorID models.ConnectorID
	Req         models.FetchNextBalancesRequest
}

func (a Activities) PluginFetchNextBalances(ctx context.Context, request FetchNextBalancesRequest) (*models.FetchNextBalancesResponse, error) {
	plugin, err := a.plugins.Get(request.ConnectorID)
	if err != nil {
		return nil, err
	}

	resp, err := plugin.FetchNextBalances(ctx, request.Req)
	if err != nil {
		// TODO(polo): temporal errors
		return nil, err
	}
	return &resp, nil
}

var PluginFetchNextBalancesActivity = Activities{}.PluginFetchNextBalances

func PluginFetchNextBalances(ctx workflow.Context, connectorID models.ConnectorID, fromPayload, state json.RawMessage, pageSize int) (*models.FetchNextBalancesResponse, error) {
	ret := models.FetchNextBalancesResponse{}
	if err := executeActivity(ctx, PluginFetchNextBalancesActivity, &ret, FetchNextBalancesRequest{
		ConnectorID: connectorID,
		Req: models.FetchNextBalancesRequest{
			FromPayload: fromPayload,
			State:       state,
			PageSize:    pageSize,
		},
	},
	); err != nil {
		return nil, err
	}
	return &ret, nil
}
