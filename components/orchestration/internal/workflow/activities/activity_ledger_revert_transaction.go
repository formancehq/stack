package activities

import (
	"context"

	sdk "github.com/formancehq/formance-sdk-go"
	"go.temporal.io/sdk/workflow"
)

type RevertTransactionRequest struct {
	Ledger string `json:"ledger"`
	TxID   int64  `json:"txId"`
}

func (a Activities) RevertTransaction(ctx context.Context, request RevertTransactionRequest) (*sdk.Transaction, error) {
	ret, _, err := a.client.TransactionsApi.
		RevertTransaction(ctx, request.Ledger, request.TxID).
		Execute()
	if err != nil {
		return nil, sdk.ExtractOpenAPIErrorMessage(err)
	}
	return &ret.Data, nil
}

var RevertTransactionActivity = Activities{}.RevertTransaction

func RevertTransaction(ctx workflow.Context, ledger string, txID int64) (*sdk.Transaction, error) {
	tx := &sdk.Transaction{}
	if err := executeActivity(ctx, RevertTransactionActivity, tx, RevertTransactionRequest{
		Ledger: ledger,
		TxID:   txID,
	}); err != nil {
		return nil, err
	}
	return tx, nil
}
