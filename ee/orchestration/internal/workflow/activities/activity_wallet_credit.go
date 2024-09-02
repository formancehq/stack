package activities

import (
	"context"
	stdtime "time"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
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
	_, err := a.client.Wallets.V1.CreditWallet(
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
			ID:             request.ID,
			IdempotencyKey: getIK(ctx),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

var CreditWalletActivity = Activities{}.CreditWallet

func CreditWallet(ctx workflow.Context, id string, request *CreditWalletRequestPayload) error {
	return executeActivity(ctx, CreditWalletActivity, nil, CreditWalletRequest{
		ID:   id,
		Data: request,
	})
}
