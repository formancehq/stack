package activities

import (
	"context"
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

var (
	ErrTransactionReferenceConflict = errors.New("TRANSACTION_REFERENCE_CONFLICT")
)

type CreateTransactionRequest struct {
	Ledger string              `json:"ledger"`
	Data   sdk.PostTransaction `json:"data"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*sdk.TransactionsResponse, error) {
	ret, httpResponse, err := a.client.TransactionsApi.
		CreateTransaction(ctx, request.Ledger).
		PostTransaction(request.Data).
		Execute()
	if err != nil {
		switch httpResponse.StatusCode {
		case http.StatusConflict:
			return nil, ErrTransactionReferenceConflict
		default:
			return nil, sdk.ExtractOpenAPIErrorMessage(err)
		}
	}
	return ret, nil
}

var CreateTransactionActivity = Activities{}.CreateTransaction

func CreateTransaction(ctx workflow.Context, ledger string, request sdk.PostTransaction) (*sdk.Transaction, error) {
	tx := &sdk.TransactionsResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data[0], nil
}
