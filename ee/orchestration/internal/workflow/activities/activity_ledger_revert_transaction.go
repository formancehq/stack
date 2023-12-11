package activities

import (
	"context"
	"math/big"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type RevertTransactionRequest struct {
	Ledger string   `json:"ledger"`
	ID     *big.Int `json:"txId"`
}

func (a Activities) RevertTransaction(ctx context.Context, request RevertTransactionRequest) (*shared.V2Transaction, error) {
	response, err := a.client.Ledger.
		V2RevertTransaction(
			ctx,
			operations.V2RevertTransactionRequest{
				Ledger: request.Ledger,
				ID:     request.ID,
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

	return &response.V2RevertTransactionResponse.Data, nil
}

var RevertTransactionActivity = Activities{}.RevertTransaction

func RevertTransaction(ctx workflow.Context, ledger string, txID *big.Int) (*shared.Transaction, error) {
	tx := &shared.Transaction{}
	if err := executeActivity(ctx, RevertTransactionActivity, tx, RevertTransactionRequest{
		Ledger: ledger,
		ID:     txID,
	}); err != nil {
		return nil, err
	}
	return tx, nil
}
