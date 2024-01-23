package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
)

func (d *DefaultClient) GetCharge(ctx context.Context, chargeID string) (*stripe.Charge, error) {
	f := connectors.ClientMetrics(ctx, "stripe", "get_charges")
	now := time.Now()
	defer f(ctx, now)

	stripe.Key = d.apiKey

	params := &stripe.ChargeParams{
		Params: stripe.Params{
			Context: ctx,
		},
	}

	chargeResponse, err := charge.Get(chargeID, params)
	if err != nil {
		return nil, err
	}

	return chargeResponse, nil
}
