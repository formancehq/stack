package activities

import (
	"context"

	sdk "github.com/formancehq/formance-sdk-go"
	"go.temporal.io/sdk/workflow"
)

type DebitWalletRequest struct {
	ID   string
	Data sdk.DebitWalletRequest
}

func (a Activities) DebitWallet(ctx context.Context, request DebitWalletRequest) (*sdk.DebitWalletResponse, error) {
	ret, _, err := a.client.WalletsApi.
		DebitWallet(ctx, request.ID).
		DebitWalletRequest(request.Data).
		Execute()
	if err != nil {
		return nil, sdk.ExtractOpenAPIErrorMessage(err)
	}
	return ret, nil
}

var DebitWalletActivity = Activities{}.DebitWallet

func DebitWallet(ctx workflow.Context, id string, request sdk.DebitWalletRequest) (*sdk.Hold, error) {
	ret := &sdk.DebitWalletResponse{}
	if err := executeActivity(ctx, DebitWalletActivity, ret, DebitWalletRequest{
		ID:   id,
		Data: request,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
