package activities

import (
	"context"
	"net/http"
	stdtime "time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
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
	ErrorResponse *shared.Error
	StatusCode    int
	RawResponse   *http.Response
}

type CreateTransactionRequest struct {
	Ledger string          `pathParam:"style=simple,explode=false,name=ledger"`
	Data   PostTransaction `request:"mediaType=application/json"`
}

type PostTransaction struct {
	Metadata  map[string]interface{}        `json:"metadata,omitempty"`
	Postings  []shared.Posting              `json:"postings,omitempty"`
	Reference *string                       `json:"reference,omitempty"`
	Script    *shared.PostTransactionScript `json:"script,omitempty"`
	Timestamp *time.Time                    `json:"timestamp,omitempty"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*shared.TransactionsResponse, error) {

	response, err := a.client.Ledger.CreateTransaction(
		ctx,
		operations.CreateTransactionRequest{
			PostTransaction: shared.PostTransaction{
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
			Ledger: request.Ledger,
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

func CreateTransaction(ctx workflow.Context, ledger string, request PostTransaction) (*shared.Transaction, error) {
	tx := &shared.TransactionsResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data[0], nil
}
