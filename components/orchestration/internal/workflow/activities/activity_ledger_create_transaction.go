package activities

import (
	"context"

	sdk "github.com/formancehq/formance-sdk-go"
	"go.temporal.io/sdk/workflow"
)

type CreateTransactionRequest struct {
	Ledger string              `json:"ledger"`
	Data   sdk.PostTransaction `json:"data"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*sdk.CreateTransactionResponse, error) {
	ret, _, err := a.client.TransactionsApi.
		CreateTransaction(ctx, request.Ledger).
		PostTransaction(request.Data).
		Execute()
	if err != nil {
		return nil, openApiErrorToApplicationError(err)
	}
	return ret, nil
}

var CreateTransactionActivity = Activities{}.CreateTransaction

func CreateTransaction(ctx workflow.Context, ledger string, request sdk.PostTransaction) (*sdk.Transaction, error) {
	tx := &sdk.CreateTransactionResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data, nil
}
