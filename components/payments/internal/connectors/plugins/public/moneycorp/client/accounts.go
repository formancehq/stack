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
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "moneycorp", "list_accounts")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/accounts", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("page[size]", strconv.Itoa(pageSize))
	q.Add("page[number]", fmt.Sprint(page))
	q.Add("sortBy", "id.asc")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			// TODO(polo): log error
			// c.logger.Error(err)
			_ = err
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return []*Account{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		// TODO(polo): retryable errors
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var accounts accountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal accounts response body: %w", err)
	}

	return accounts.Accounts, nil
}
