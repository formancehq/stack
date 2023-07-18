package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Account struct {
	ID          string    `json:"id"`
	AccountName string    `json:"account_name"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Client) GetAccounts(ctx context.Context, page int) ([]*Account, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.buildEndpoint("v2/accounts/find?page=%d&per_page=25", page), http.NoBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//nolint:tagliatelle // allow for client code
	type response struct {
		Accounts   []*Account `json:"accounts"`
		Pagination struct {
			NextPage int `json:"next_page"`
		} `json:"pagination"`
	}

	var res response
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, 0, err
	}

	return res.Accounts, res.Pagination.NextPage, nil
}
