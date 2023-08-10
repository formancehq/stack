package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type accountsResponse struct {
	Accounts []*Account `json:"data"`
}

type Account struct {
	ID         string `json:"id"`
	Attributes struct {
		AccountName string `json:"accountName"`
	} `json:"attributes"`
}

func (c *Client) GetAccounts(ctx context.Context, page int, pageSize int) ([]*Account, error) {
	endpoint := fmt.Sprintf("%s/accounts", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("pagesize", strconv.Itoa(pageSize))
	q.Add("pagenumber", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	var accounts accountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal accounts response body: %w", err)
	}

	return accounts.Accounts, nil
}
