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
	externalAccountsEndpoint = "https://api.stripe.com/v1/accounts/%s/external_accounts"
)

func (d *DefaultClient) ExternalAccounts(ctx context.Context, options ...ClientOption) ([]*stripe.ExternalAccount, bool, error) {
	if d.stripeAccount == "" {
		return nil, false, errors.New("stripe account is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(externalAccountsEndpoint, d.stripeAccount), nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "creating http request")
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
		return nil, false, errors.Wrap(err, "doing request")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	type listResponse struct {
		TransactionsListResponse
		Data []json.RawMessage `json:"data"`
	}

	rsp := &listResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(rsp)
	if err != nil {
		return nil, false, errors.Wrap(err, "decoding response")
	}

	externalAccounts := make([]*stripe.ExternalAccount, 0)

	if len(rsp.Data) > 0 {
		for _, data := range rsp.Data {
			account := &stripe.ExternalAccount{}

			err = json.Unmarshal(data, &account)
			if err != nil {
				return nil, false, err
			}

			externalAccounts = append(externalAccounts, account)
		}
	}

	return externalAccounts, rsp.HasMore, nil
}
