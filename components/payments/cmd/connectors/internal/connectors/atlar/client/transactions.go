package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	transactionsEndpoint = "https://api.atlar.com/v1/transactions"
)

type Transaction struct {
	Id string `json:"id"`
}

type TransactionsListResponse struct {
	ListResponse
	Items []*Transaction `json:"items"`
}

func (d *DefaultClient) Transactions(ctx context.Context,
	options ...ClientOption,
) ([]*Transaction, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, transactionsEndpoint, nil)
	if err != nil {
		return nil, "", errors.Wrap(err, "creating http request")
	}

	for _, opt := range options {
		opt.Apply(req)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(d.accessKey, d.secret)

	var httpResponse *http.Response

	httpResponse, err = d.httpClient.Do(req)
	if err != nil {
		return nil, "", errors.Wrap(err, "doing request")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	type listResponse struct {
		TransactionsListResponse
		Items []json.RawMessage `json:"items"`
	}

	rsp := &listResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(rsp)
	if err != nil {
		return nil, "", errors.Wrap(err, "decoding response")
	}

	transactions := make([]*Transaction, 0)

	if len(rsp.Items) > 0 {
		for _, item := range rsp.Items {
			transaction := &Transaction{}

			err = json.Unmarshal(item, &transaction)
			if err != nil {
				return nil, "", err
			}

			transactions = append(transactions, transaction)
		}
	}

	return transactions, rsp.NextToken, nil
}
