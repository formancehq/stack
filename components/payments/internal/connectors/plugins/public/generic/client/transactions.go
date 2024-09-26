package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/genericclient"
)

func (c *Client) ListTransactions(ctx context.Context, page, pageSize int64, updatedAtFrom time.Time) ([]genericclient.Transaction, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "generic", "list_transactions")
	// now := time.Now()
	// defer f(ctx, now)

	req := c.apiClient.DefaultApi.GetTransactions(ctx).
		Page(page).
		PageSize(pageSize)

	if !updatedAtFrom.IsZero() {
		req = req.UpdatedAtFrom(updatedAtFrom)
	}

	transactions, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
