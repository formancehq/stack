package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/genericclient"
)

func (c *Client) ListAccounts(ctx context.Context, page, pageSize int64, createdAtFrom time.Time) ([]genericclient.Account, error) {
	// TODO(f): Add metrics
	// f := connectors.ClientMetrics(ctx, "generic", "list_accounts")
	// now := time.Now()
	// defer f(ctx, now)

	req := c.apiClient.DefaultApi.
		GetAccounts(ctx).
		Page(page).
		PageSize(pageSize)

	if !createdAtFrom.IsZero() {
		req = req.CreatedAtFrom(createdAtFrom)
	}

	accounts, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
