package activities

import (
	"context"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type GetAccountRequest struct {
	Ledger string `json:"ledger"`
	ID     string `json:"id"`
}

func (a Activities) GetAccount(ctx context.Context, request GetAccountRequest) (*shared.V2AccountResponse, error) {
	response, err := a.client.Ledger.V2GetAccount(
		ctx,
		operations.V2GetAccountRequest{
			Address: request.ID,
			Ledger:  request.Ledger,
		},
	)
	if err != nil {
		switch err := err.(type) {
		case *sdkerrors.V2ErrorResponse:
			return nil, temporal.NewApplicationError(err.ErrorMessage, string(err.ErrorCode), err.Details)
		default:
			return nil, err
		}
	}

	return response.V2AccountResponse, nil
}

var GetAccountActivity = Activities{}.GetAccount

func GetAccount(ctx workflow.Context, ledger, id string) (*shared.V2Account, error) {
	ret := &shared.V2AccountResponse{}
	if err := executeActivity(ctx, GetAccountActivity, ret, GetAccountRequest{
		Ledger: ledger,
		ID:     id,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
