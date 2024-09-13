package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
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
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "moneycorp", "list_account_balances")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/accounts/%s/balances", c.endpoint, accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	balances := balancesResponse{Balances: make([]*Balance, 0)}
	var errRes moneycorpError

	_, err = c.httpClient.Do(req, &balances, &errRes)
	switch err {
	case nil:
		return balances.Balances, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to get account balances: %w", err)
}
