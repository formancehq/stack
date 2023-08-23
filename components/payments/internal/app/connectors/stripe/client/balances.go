package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
)

const (
	balanceEndpoint = "https://api.stripe.com/v1/balance"
)

type BalanceResponse struct {
	*stripe.Balance
}

func (d *DefaultClient) Balance(ctx context.Context, options ...ClientOption) (*stripe.Balance, error) {
	if d.stripeAccount == "" {
		return nil, errors.New("stripe account is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, balanceEndpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	for _, opt := range options {
		opt.Apply(req)
	}

	req.Header.Set("Stripe-Account", d.stripeAccount)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(d.apiKey, "") // gfyrag: really weird authentication right?

	var httpResponse *http.Response
	httpResponse, err = d.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	type balanceResponse struct {
		BalanceResponse
	}

	rsp := &balanceResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(rsp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return rsp.Balance, nil
}
