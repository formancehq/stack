package activities

import (
	"context"
	"fmt"
	"net/http"
	stdtime "time"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type DebitWalletRequest struct {
	ID   string                     `json:"id"`
	Data *DebitWalletRequestPayload `json:"data"`
}

type DebitWalletRequestPayload struct {
	Amount      shared.Monetary `json:"amount"`
	Balances    []string        `json:"balances,omitempty"`
	Description *string         `json:"description,omitempty"`
	Destination *shared.Subject `json:"destination,omitempty"`
	// Metadata associated with the wallet.
	Metadata map[string]string `json:"metadata"`
	// Set to true to create a pending hold. If false, the wallet will be debited immediately.
	Pending *bool `json:"pending,omitempty"`
	// cannot be used in conjunction with `pending` property
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

func (a Activities) DebitWallet(ctx context.Context, request DebitWalletRequest) (*shared.DebitWalletResponse, error) {
	response, err := a.client.Wallets.DebitWallet(
		ctx,
		operations.DebitWalletRequest{
			DebitWalletRequest: &shared.DebitWalletRequest{
				Amount:      request.Data.Amount,
				Balances:    request.Data.Balances,
				Description: request.Data.Description,
				Destination: request.Data.Destination,
				Metadata:    request.Data.Metadata,
				Pending:     request.Data.Pending,
				Timestamp: func() *stdtime.Time {
					if request.Data.Timestamp == nil {
						return nil
					}
					return &request.Data.Timestamp.Time
				}(),
			},
			ID: request.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	if response.WalletsErrorResponse != nil {
		return nil, temporal.NewApplicationError(
			response.WalletsErrorResponse.ErrorMessage,
			string(response.WalletsErrorResponse.ErrorCode),
		)
	}

	switch response.StatusCode {
	case http.StatusNoContent, http.StatusCreated:
		return response.DebitWalletResponse, nil
	default:
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var DebitWalletActivity = Activities{}.DebitWallet

func DebitWallet(ctx workflow.Context, id string, request *DebitWalletRequestPayload) (*shared.Hold, error) {
	ret := &shared.DebitWalletResponse{}
	if err := executeActivity(ctx, DebitWalletActivity, ret, DebitWalletRequest{
		ID:   id,
		Data: request,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
