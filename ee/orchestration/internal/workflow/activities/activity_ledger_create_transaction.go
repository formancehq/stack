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

// CreateTransactionResponse - OK
type CreateTransactionResponse struct {
	Data []shared.Transaction `json:"data"`
}

type CreateTransactionWrapper struct {
	ContentType string
	// OK
	CreateTransactionResponse *CreateTransactionResponse
	// Error
	ErrorResponse *shared.ErrorResponse
	StatusCode    int
	RawResponse   *http.Response
}

type CreateTransactionRequest struct {
	Ledger string                   `pathParam:"style=simple,explode=false,name=ledger"`
	Data   shared.V2PostTransaction `request:"mediaType=application/json"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*shared.V2CreateTransactionResponse, error) {

	response, err := a.client.Ledger.V2CreateTransaction(
		ctx,
		operations.V2CreateTransactionRequest{
			V2PostTransaction: request.Data,
			Ledger:            request.Ledger,
		},
	)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return response.V2CreateTransactionResponse, nil
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

var CreateTransactionActivity = Activities{}.CreateTransaction

func CreateTransaction(ctx workflow.Context, ledger string, request shared.V2PostTransaction) (*shared.V2Transaction, error) {
	tx := &shared.V2CreateTransactionResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data, nil
}
