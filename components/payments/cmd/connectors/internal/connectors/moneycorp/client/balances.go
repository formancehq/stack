package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type balancesResponse struct {
	Balances []*Balance `json:"data"`
}

type Balance struct {
	ID         string `json:"id"`
	Attributes struct {
		CurrencyCode     string      `json:"currencyCode"`
		OverallBalance   json.Number `json:"overallBalance"`
		AvailableBalance json.Number `json:"availableBalance"`
		ClearedBalance   json.Number `json:"clearedBalance"`
		ReservedBalance  json.Number `json:"reservedBalance"`
		UnclearedBalance json.Number `json:"unclearedBalance"`
	} `json:"attributes"`
}

func (c *Client) GetAccountBalances(ctx context.Context, accountID string) ([]*Balance, error) {
	f := connectors.ClientMetrics(ctx, "moneycorp", "list_account_balances")
	now := time.Now()
	defer f(ctx, now)

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
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var balances balancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return nil, fmt.Errorf("failed to unmarshal balances response body: %w", err)
	}

	return balances.Balances, nil
}
