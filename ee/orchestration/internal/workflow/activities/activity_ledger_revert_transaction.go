package activities

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
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
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusCreated:
		return &response.V2RevertTransactionResponse.Data, nil
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
