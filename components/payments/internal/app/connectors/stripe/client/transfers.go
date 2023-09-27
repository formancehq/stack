package client

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/transfer"
)

type CreateTransferRequest struct {
	Amount      int64
	Currency    string
	Destination string
	Description string
}

func (d *DefaultClient) CreateTransfer(ctx context.Context, createTransferRequest *CreateTransferRequest, options ...ClientOption) (*stripe.Transfer, error) {
	stripe.Key = d.apiKey

	params := &stripe.TransferParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Amount:      stripe.Int64(createTransferRequest.Amount),
		Currency:    stripe.String(createTransferRequest.Currency),
		Destination: stripe.String(createTransferRequest.Destination),
	}

	if d.stripeAccount != "" {
		params.SetStripeAccount(d.stripeAccount)
	}

	transferResponse, err := transfer.New(params)
	if err != nil {
		return nil, errors.Wrap(err, "creating transfer")
	}

	return transferResponse, nil
}
