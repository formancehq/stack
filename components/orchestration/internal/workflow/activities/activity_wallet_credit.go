package activities

import (
	"context"
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

type CreditWalletRequest struct {
	ID   string                  `json:"id"`
	Data sdk.CreditWalletRequest `json:"data"`
}

func (a Activities) CreditWallet(ctx context.Context, request CreditWalletRequest) error {
	httpResponse, err := a.client.WalletsApi.
		CreditWallet(ctx, request.ID).
		CreditWalletRequest(request.Data).
		Execute()
	if err != nil {
		switch httpResponse.StatusCode {
		case http.StatusNotFound:
			return errors.New("wallet not found")
		default:
			return openApiErrorToApplicationError(err)
		}
	}
	return nil
}

var CreditWalletActivity = Activities{}.CreditWallet

func CreditWallet(ctx workflow.Context, id string, request sdk.CreditWalletRequest) error {
	return executeActivity(ctx, CreditWalletActivity, nil, CreditWalletRequest{
		ID:   id,
		Data: request,
	})
}
