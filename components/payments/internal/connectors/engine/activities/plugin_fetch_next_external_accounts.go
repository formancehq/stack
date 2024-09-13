package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

type FetchNextExternalAccountsRequest struct {
	ConnectorID models.ConnectorID
	Req         models.FetchNextExternalAccountsRequest
}

func (a Activities) PluginFetchNextExternalAccounts(ctx context.Context, request FetchNextExternalAccountsRequest) (*models.FetchNextExternalAccountsResponse, error) {
	plugin, err := a.plugins.Get(request.ConnectorID)
	if err != nil {
		return nil, err
	}

	resp, err := plugin.FetchNextExternalAccounts(ctx, request.Req)
	if err != nil {
		// TODO(polo): temporal errors
		return nil, err
	}

	return &resp, nil
}

var PluginFetchNextExternalAccountsActivity = Activities{}.PluginFetchNextExternalAccounts

func PluginFetchNextExternalAccounts(ctx workflow.Context, connectorID models.ConnectorID, fromPayload, state json.RawMessage, pageSize int) (*models.FetchNextExternalAccountsResponse, error) {
	ret := models.FetchNextExternalAccountsResponse{}
	if err := executeActivity(ctx, PluginFetchNextExternalAccountsActivity, &ret, FetchNextExternalAccountsRequest{
		ConnectorID: connectorID,
		Req: models.FetchNextExternalAccountsRequest{
			FromPayload: fromPayload,
			State:       state,
			PageSize:    pageSize,
		},
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}
