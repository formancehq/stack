package client

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/payout"
)

type CreatePayoutRequest struct {
	IdempotencyKey string
	Amount         int64
	Currency       string
	Destination    string
	Description    string
}

func (d *DefaultClient) CreatePayout(ctx context.Context, createPayoutRequest *CreatePayoutRequest, options ...ClientOption) (*stripe.Payout, error) {
	stripe.Key = d.apiKey

	params := &stripe.PayoutParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Amount:      stripe.Int64(createPayoutRequest.Amount),
		Currency:    stripe.String(createPayoutRequest.Currency),
		Destination: stripe.String(createPayoutRequest.Destination),
		Method:      stripe.String("standard"),
	}

	if d.stripeAccount != "" {
		params.SetStripeAccount(d.stripeAccount)
	}

	if createPayoutRequest.IdempotencyKey != "" {
		params.IdempotencyKey = stripe.String(createPayoutRequest.IdempotencyKey)
	}

	if createPayoutRequest.Description != "" {
		params.Description = stripe.String(createPayoutRequest.Description)
	}

	payoutResponse, err := payout.New(params)
	if err != nil {
		return nil, errors.Wrap(err, "creating transfer")
	}

	return payoutResponse, nil
}

func (d *DefaultClient) GetPayout(ctx context.Context, payoutID string, options ...ClientOption) (*stripe.Payout, error) {
	stripe.Key = d.apiKey

	payoutResponse, err := payout.Get(payoutID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting payout")
	}

	return payoutResponse, nil
}
