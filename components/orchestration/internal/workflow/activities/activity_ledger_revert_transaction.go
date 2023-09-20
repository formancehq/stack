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

type RevertTransactionRequest struct {
	Ledger string `json:"ledger"`
	ID     int64  `json:"txId"`
}

func (a Activities) RevertTransaction(ctx context.Context, request RevertTransactionRequest) (*shared.Transaction, error) {
	response, err := a.client.Ledger.
		RevertTransaction(
			ctx,
			operations.RevertTransactionRequest{
				Ledger: request.Ledger,
				ID:     request.ID,
			},
		)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusCreated:
		return &response.RevertTransactionResponse.Data, nil
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

var RevertTransactionActivity = Activities{}.RevertTransaction

func RevertTransaction(ctx workflow.Context, ledger string, txID int64) (*shared.Transaction, error) {
	tx := &shared.Transaction{}
	if err := executeActivity(ctx, RevertTransactionActivity, tx, RevertTransactionRequest{
		Ledger: ledger,
		ID:     txID,
	}); err != nil {
		return nil, err
	}
	return tx, nil
}
