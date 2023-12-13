package activities

import (
	"context"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"

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
	Ledger string                 `pathParam:"style=simple,explode=false,name=ledger"`
	Data   shared.PostTransaction `request:"mediaType=application/json"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*shared.TransactionsResponse, error) {

	response, err := a.client.Ledger.CreateTransaction(
		ctx,
		operations.CreateTransactionRequest{
			PostTransaction: request.Data,
			Ledger:          request.Ledger,
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

	return response.TransactionsResponse, nil
}

var CreateTransactionActivity = Activities{}.CreateTransaction

func CreateTransaction(ctx workflow.Context, ledger string, request shared.PostTransaction) (*shared.Transaction, error) {
	tx := &shared.TransactionsResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data[0], nil
}
