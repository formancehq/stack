package activities

import (
	"context"
	"errors"
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go"
	"go.temporal.io/sdk/workflow"
)

type GetAccountRequest struct {
	Ledger string `json:"ledger"`
	ID     string `json:"id"`
}

func (a Activities) GetAccount(ctx context.Context, request GetAccountRequest) (*sdk.AccountResponse, error) {
	ret, httpResponse, err := a.client.AccountsApi.
		GetAccount(ctx, request.Ledger, request.ID).
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

var GetAccountActivity = Activities{}.GetAccount

func GetAccount(ctx workflow.Context, ledger, id string) (*sdk.AccountWithVolumesAndBalances, error) {
	ret := &sdk.AccountResponse{}
	if err := executeActivity(ctx, GetAccountActivity, ret, GetAccountRequest{
		Ledger: ledger,
		ID:     id,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
