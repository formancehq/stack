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
	balanceTransactionsEndpoint = "https://api.stripe.com/v1/balance_transactions"
)

//nolint:tagliatelle // allow different styled tags in client
type TransactionsListResponse struct {
	HasMore bool                         `json:"has_more"`
	Data    []*stripe.BalanceTransaction `json:"data"`
}

func (d *DefaultClient) BalanceTransactions(ctx context.Context,
	options ...ClientOption,
) ([]*stripe.BalanceTransaction, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, balanceTransactionsEndpoint, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "creating http request")
	}

	for _, opt := range options {
		opt.Apply(req)
	}

	if d.stripeAccount != "" {
		req.Header.Set("Stripe-Account", d.stripeAccount)
	}

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

	asBalanceTransactions := make([]*stripe.BalanceTransaction, 0)

	if len(rsp.Data) > 0 {
		for _, data := range rsp.Data {
			asBalanceTransaction := &stripe.BalanceTransaction{}

			err = json.Unmarshal(data, &asBalanceTransaction)
			if err != nil {
				return nil, false, err
			}

			asBalanceTransactions = append(asBalanceTransactions, asBalanceTransaction)
		}
	}

	return asBalanceTransactions, rsp.HasMore, nil
}
