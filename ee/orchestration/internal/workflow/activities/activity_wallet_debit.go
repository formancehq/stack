package activities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type DebitWalletRequest struct {
	ID   string                     `json:"id"`
	Data *shared.DebitWalletRequest `json:"data"`
}

func (a Activities) DebitWallet(ctx context.Context, request DebitWalletRequest) (*shared.DebitWalletResponse, error) {
	response, err := a.client.Wallets.DebitWallet(
		ctx,
		operations.DebitWalletRequest{
			DebitWalletRequest: request.Data,
			ID:                 request.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	if response.WalletsErrorResponse != nil {
		return nil, temporal.NewApplicationError(
			response.WalletsErrorResponse.ErrorMessage,
			string(response.WalletsErrorResponse.ErrorCode),
		)
	}

	switch response.StatusCode {
	case http.StatusNoContent, http.StatusCreated:
		return response.DebitWalletResponse, nil
	default:
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var DebitWalletActivity = Activities{}.DebitWallet

func DebitWallet(ctx workflow.Context, id string, request *shared.DebitWalletRequest) (*shared.Hold, error) {
	ret := &shared.DebitWalletResponse{}
	if err := executeActivity(ctx, DebitWalletActivity, ret, DebitWalletRequest{
		ID:   id,
		Data: request,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
