package activities

import (
	"context"

	"github.com/formancehq/go-libs/pointer"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"go.temporal.io/sdk/workflow"
)

type ListWalletsRequest struct {
	Name string `json:"name"`
}

func (a Activities) ListWallets(ctx context.Context, request ListWalletsRequest) (*shared.ListWalletsResponse, error) {
	response, err := a.client.Wallets.V1.ListWallets(
		ctx,
		operations.ListWalletsRequest{
			Name: pointer.For(request.Name),
		},
	)
	if err != nil {
		return nil, err
	}

	return response.ListWalletsResponse, nil
}

var ListWalletsActivity = Activities{}.ListWallets

func ListWallets(ctx workflow.Context, request ListWalletsRequest) (*shared.ListWalletsResponse, error) {
	ret := &shared.ListWalletsResponse{}
	if err := executeActivity(ctx, ListWalletsActivity, ret, request); err != nil {
		return nil, err
	}
	return ret, nil
}
