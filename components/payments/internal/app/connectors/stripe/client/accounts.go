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
	accountsEndpoint = "https://api.stripe.com/v1/accounts"
)

//nolint:tagliatelle // allow different styled tags in client
type AccountsListResponse struct {
	HasMore bool              `json:"has_more"`
	Data    []*stripe.Account `json:"data"`
}

func (d *DefaultClient) Accounts(ctx context.Context,
	options ...ClientOption,
) ([]*stripe.Account, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accountsEndpoint, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "creating http request")
	}

	for _, opt := range options {
		opt.Apply(req)
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
		AccountsListResponse
		Data []json.RawMessage `json:"data"`
	}

	rsp := &listResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(rsp)
	if err != nil {
		return nil, false, errors.Wrap(err, "decoding response")
	}

	accounts := make([]*stripe.Account, 0)

	if len(rsp.Data) > 0 {
		for _, data := range rsp.Data {
			account := &stripe.Account{}

			err = json.Unmarshal(data, &account)
			if err != nil {
				return nil, false, err
			}

			accounts = append(accounts, account)
		}
	}

	return accounts, rsp.HasMore, nil
}
