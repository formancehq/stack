package activities

import (
	"context"
	stdtime "time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/go-libs/time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// CreateTransactionResponse - OK
type CreateTransactionResponse struct {
	Data []shared.Transaction `json:"data"`
}

type CreateTransactionRequest struct {
	Ledger string          `pathParam:"style=simple,explode=false,name=ledger"`
	Data   PostTransaction `request:"mediaType=application/json"`
}

type PostTransaction struct {
	Metadata  map[string]string               `json:"metadata,omitempty"`
	Postings  []shared.V2Posting              `json:"postings,omitempty"`
	Reference *string                         `json:"reference,omitempty"`
	Script    *shared.V2PostTransactionScript `json:"script,omitempty"`
	Timestamp *time.Time                      `json:"timestamp,omitempty"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*shared.V2CreateTransactionResponse, error) {

	response, err := a.client.Ledger.V2.CreateTransaction(
		ctx,
		operations.V2CreateTransactionRequest{
			V2PostTransaction: shared.V2PostTransaction{
				Metadata:  request.Data.Metadata,
				Postings:  request.Data.Postings,
				Reference: request.Data.Reference,
				Script:    request.Data.Script,
				Timestamp: func() *stdtime.Time {
					if request.Data.Timestamp == nil {
						return nil
					}
					return &request.Data.Timestamp.Time
				}(),
			},
			Ledger:         request.Ledger,
			IdempotencyKey: getIK(ctx),
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

	return response.V2CreateTransactionResponse, nil
}

var CreateTransactionActivity = Activities{}.CreateTransaction

func CreateTransaction(ctx workflow.Context, ledger string, request PostTransaction) (*shared.V2Transaction, error) {
	tx := &shared.V2CreateTransactionResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data, nil
}
