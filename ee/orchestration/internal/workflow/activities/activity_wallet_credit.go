package activities

import (
	"context"
	"fmt"
	"net/http"
	stdtime "time"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type CreditWalletRequest struct {
	ID   string                      `json:"id"`
	Data *CreditWalletRequestPayload `json:"data"`
}

type CreditWalletRequestPayload struct {
	Amount shared.Monetary `json:"amount"`
	// The balance to credit
	Balance *string `json:"balance,omitempty"`
	// Metadata associated with the wallet.
	Metadata  map[string]string `json:"metadata"`
	Reference *string           `json:"reference,omitempty"`
	Sources   []shared.Subject  `json:"sources"`
	Timestamp *time.Time        `json:"timestamp,omitempty"`
}

func (a Activities) CreditWallet(ctx context.Context, request CreditWalletRequest) error {
	response, err := a.client.Wallets.CreditWallet(
		ctx,
		operations.CreditWalletRequest{
			CreditWalletRequest: &shared.CreditWalletRequest{
				Amount:    request.Data.Amount,
				Balance:   request.Data.Balance,
				Metadata:  request.Data.Metadata,
				Reference: request.Data.Reference,
				Sources:   request.Data.Sources,
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
		return err
	}

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return errors.New("wallet not found")
	default:
		if response.WalletsErrorResponse != nil {
			return temporal.NewApplicationError(
				response.WalletsErrorResponse.ErrorMessage,
				string(response.WalletsErrorResponse.ErrorCode),
			)
		}

		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var CreditWalletActivity = Activities{}.CreditWallet

func CreditWallet(ctx workflow.Context, id string, request *CreditWalletRequestPayload) error {
	return executeActivity(ctx, CreditWalletActivity, nil, CreditWalletRequest{
		ID:   id,
		Data: request,
	})
}
