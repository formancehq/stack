package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/genericclient"
)

func (c *Client) GetBalances(ctx context.Context, accountID string) (*genericclient.Balances, error) {
	f := connectors.ClientMetrics(ctx, "generic", "get_balance")
	now := time.Now()
	defer f(ctx, now)

	req := c.apiClient.DefaultApi.GetAccountBalances(ctx, accountID)

	balances, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	return balances, nil
}
