package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type balancesResponse struct {
	Balances []*Balance `json:"data"`
}

type Balance struct {
	ID         string `json:"id"`
	Attributes struct {
		CurrencyCode     string  `json:"currencyCode"`
		OverallBalance   float64 `json:"overallBalance"`
		AvailableBalance float64 `json:"availableBalance"`
		ClearedBalance   float64 `json:"clearedBalance"`
		ReservedBalance  float64 `json:"reservedBalance"`
		UnclearedBalance float64 `json:"unclearedBalance"`
	} `json:"attributes"`
}

func (c *Client) GetAccountBalances(ctx context.Context, accountID string) ([]*Balance, error) {
	endpoint := fmt.Sprintf("%s/accounts/%s/balances", c.endpoint, accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balances: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get balances: %s", resp.Status)
	}

	var balances balancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return nil, fmt.Errorf("failed to unmarshal balances response body: %w", err)
	}

	return balances.Balances, nil
}
