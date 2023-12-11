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
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return response.V2AccountResponse, nil
	case http.StatusNotFound:
		return nil, errors.New("wallet not found")
	default:
		if response.V2ErrorResponse != nil {
			return nil, temporal.NewApplicationError(
				response.V2ErrorResponse.ErrorMessage,
				string(response.V2ErrorResponse.ErrorCode),
				response.V2ErrorResponse.Details)
		}

		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
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
