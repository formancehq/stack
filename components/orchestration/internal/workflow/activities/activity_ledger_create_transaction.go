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

type CreateTransactionRequest struct {
	Ledger string                 `json:"ledger"`
	Data   shared.PostTransaction `json:"data"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*shared.CreateTransactionResponse, error) {
	response, err := a.client.Transactions.
		CreateTransaction(
			ctx,
			operations.CreateTransactionRequest{
				PostTransaction: request.Data,
				Ledger:          request.Ledger,
			},
		)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return response.CreateTransactionResponse, nil
	default:
		if response.ErrorResponse != nil {
			return nil, temporal.NewApplicationError(
				response.ErrorResponse.ErrorMessage,
				string(response.ErrorResponse.ErrorCode),
				response.ErrorResponse.Details)
		}

		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var CreateTransactionActivity = Activities{}.CreateTransaction

func CreateTransaction(ctx workflow.Context, ledger string, request shared.PostTransaction) (*shared.Transaction, error) {
	tx := &shared.CreateTransactionResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data, nil
}
