package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	pageSize = 100
)

type accountsResponse struct {
	Accounts []account `json:"data"`
}

type account struct {
	ID string `json:"id"`
}

func (c *Client) GetAllAccounts(ctx context.Context) ([]account, error) {
	var accounts []account

	for page := 0; ; page++ {
		pagedAccounts, err := c.getAccounts(ctx, page, pageSize)
		if err != nil {
			return nil, err
		}

		if len(pagedAccounts) == 0 {
			break
		}

		accounts = append(accounts, pagedAccounts...)

		if len(pagedAccounts) < pageSize {
			break
		}
	}

	return accounts, nil
}

func (c *Client) getAccounts(ctx context.Context, page int, pageSize int) ([]account, error) {
	endpoint := fmt.Sprintf("%s/accounts", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("pagesize", strconv.Itoa(pageSize))
	q.Add("pagenumber", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	var accounts accountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal login response body: %w", err)
	}

	return accounts.Accounts, nil
}
