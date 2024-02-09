package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Balance struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Currency  string    `json:"currency"`
	Amount    string    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Client) GetBalances(ctx context.Context, page int) ([]*Balance, int, error) {
	f := connectors.ClientMetrics(ctx, "currencycloud", "list_balances")
	now := time.Now()
	defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.buildEndpoint("v2/balances/find"), http.NoBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page", fmt.Sprint(page))
	q.Add("per_page", "25")
	q.Add("order", "created_at")
	q.Add("order_asc_desc", "asc")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	//nolint:tagliatelle // allow for client code
	type response struct {
		Balances   []*Balance `json:"balances"`
		Pagination struct {
			NextPage int `json:"next_page"`
		} `json:"pagination"`
	}

	var res response
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, 0, err
	}

	return res.Balances, res.Pagination.NextPage, nil
}
