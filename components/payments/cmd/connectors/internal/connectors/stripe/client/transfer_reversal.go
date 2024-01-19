package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/reversal"
)

type CreateTransferReversalRequest struct {
	TransferID  string
	Amount      int64
	Description string
	Metadata    map[string]string
}

func (d *DefaultClient) ReverseTransfer(ctx context.Context, createTransferReversalRequest *CreateTransferReversalRequest, options ...ClientOption) (*stripe.Reversal, error) {
	f := connectors.ClientMetrics(ctx, "stripe", "reverse_transfer")
	now := time.Now()
	defer f(ctx, now)

	stripe.Key = d.apiKey

	params := &stripe.ReversalParams{
		Params: stripe.Params{
			Context:  ctx,
			Metadata: createTransferReversalRequest.Metadata,
		},
		Transfer:    stripe.String(createTransferReversalRequest.TransferID),
		Amount:      stripe.Int64(createTransferReversalRequest.Amount),
		Description: stripe.String(createTransferReversalRequest.Description),
	}

	if d.stripeAccount != "" {
		params.SetStripeAccount(d.stripeAccount)
	}

	reversalResponse, err := reversal.New(params)
	if err != nil {
		return nil, err
	}

	return reversalResponse, nil
}
