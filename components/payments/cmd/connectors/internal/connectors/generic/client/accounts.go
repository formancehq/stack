package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/genericclient"
)

func (c *Client) ListAccounts(ctx context.Context, page, pageSize int64) ([]genericclient.Account, error) {
	f := connectors.ClientMetrics(ctx, "generic", "list_accounts")
	now := time.Now()
	defer f(ctx, now)

	req := c.apiClient.DefaultApi.
		GetAccounts(ctx).
		Page(page).
		PageSize(pageSize)

	accounts, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
