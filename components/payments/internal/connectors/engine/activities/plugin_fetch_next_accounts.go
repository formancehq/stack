package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

type FetchNextAccountsRequest struct {
	ConnectorID models.ConnectorID
	Req         models.FetchNextAccountsRequest
}

func (a Activities) PluginFetchNextAccounts(ctx context.Context, request FetchNextAccountsRequest) (*models.FetchNextAccountsResponse, error) {
	plugin, err := a.plugins.Get(request.ConnectorID)
	if err != nil {
		return nil, err
	}

	resp, err := plugin.FetchNextAccounts(ctx, request.Req)
	if err != nil {
		// TODO(polo): temporal errors
		return nil, err
	}
	return &resp, nil
}

var PluginFetchNextAccountsActivity = Activities{}.PluginFetchNextAccounts

func PluginFetchNextAccounts(ctx workflow.Context, connectorID models.ConnectorID, fromPayload, state json.RawMessage, pageSize int) (*models.FetchNextAccountsResponse, error) {
	ret := models.FetchNextAccountsResponse{}
	if err := executeActivity(ctx, PluginFetchNextAccountsActivity, &ret, FetchNextAccountsRequest{
		ConnectorID: connectorID,
		Req: models.FetchNextAccountsRequest{
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
