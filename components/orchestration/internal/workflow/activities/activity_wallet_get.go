package activities

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type GetWalletRequest struct {
	ID string `json:"id"`
}

func (a Activities) GetWallet(ctx context.Context, request GetWalletRequest) (*shared.GetWalletResponse, error) {
	response, err := a.client.Wallets.GetWallet(
		ctx,
		operations.GetWalletRequest{
			ID: request.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return response.GetWalletResponse, nil
	case http.StatusNotFound:
		return nil, errors.New("wallet not found")
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

var GetWalletActivity = Activities{}.GetWallet

func GetWallet(ctx workflow.Context, id string) (*shared.WalletWithBalances, error) {
	ret := &shared.GetWalletResponse{}
	if err := executeActivity(ctx, GetWalletActivity, ret, GetWalletRequest{
		ID: id,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
