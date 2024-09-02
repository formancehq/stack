package activities

import (
	"context"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type GetAccountRequest struct {
	Ledger string `json:"ledger"`
	ID     string `json:"id"`
}

func (a Activities) GetAccount(ctx context.Context, request GetAccountRequest) (*shared.AccountResponse, error) {
	response, err := a.client.Ledger.V1.GetAccount(
		ctx,
		operations.GetAccountRequest{
			Address: request.ID,
			Ledger:  request.Ledger,
		},
	)
	if err != nil {
		switch err := err.(type) {
		case *sdkerrors.ErrorResponse:
			return nil, temporal.NewApplicationError(err.ErrorMessage, string(err.ErrorCode), err.Details)
		default:
			return nil, err
		}
	}

	return response.AccountResponse, nil
}

var GetAccountActivity = Activities{}.GetAccount

func GetAccount(ctx workflow.Context, ledger, id string) (*shared.AccountWithVolumesAndBalances, error) {
	ret := &shared.AccountResponse{}
	if err := executeActivity(ctx, GetAccountActivity, ret, GetAccountRequest{
		Ledger: ledger,
		ID:     id,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
