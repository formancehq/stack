package activities

import (
	"context"
	"errors"
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go"
	"go.temporal.io/sdk/workflow"
)

type GetWalletRequest struct {
	ID string `json:"id"`
}

func (a Activities) GetWallet(ctx context.Context, request GetWalletRequest) (*sdk.GetWalletResponse, error) {
	ret, httpResponse, err := a.client.WalletsApi.
		GetWallet(ctx, request.ID).
		Execute()
	if err != nil {
		switch httpResponse.StatusCode {
		case http.StatusNotFound:
			return nil, errors.New("wallet not found")
		default:
			return nil, openApiErrorToApplicationError(err)
		}
	}
	return ret, nil
}

var GetWalletActivity = Activities{}.GetWallet

func GetWallet(ctx workflow.Context, id string) (*sdk.WalletWithBalances, error) {
	ret := &sdk.GetWalletResponse{}
	if err := executeActivity(ctx, GetWalletActivity, ret, GetWalletRequest{
		ID: id,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
