package activities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ListWalletsRequest struct {
	Name string `json:"name"`
}

func (a Activities) ListWallets(ctx context.Context, request ListWalletsRequest) (*shared.ListWalletsResponse, error) {
	response, err := a.client.Wallets.ListWallets(
		ctx,
		operations.ListWalletsRequest{
			Name: pointer.For(request.Name),
		},
	)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return response.ListWalletsResponse, nil
	default:
		if response.WalletsErrorResponse != nil {
			return nil, temporal.NewApplicationError(
				response.WalletsErrorResponse.ErrorMessage,
				string(response.WalletsErrorResponse.ErrorCode),
			)
		}

		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var ListWalletsActivity = Activities{}.ListWallets

func ListWallets(ctx workflow.Context, request ListWalletsRequest) (*shared.ListWalletsResponse, error) {
	ret := &shared.ListWalletsResponse{}
	if err := executeActivity(ctx, ListWalletsActivity, ret, request); err != nil {
		return nil, err
	}
	return ret, nil
}
