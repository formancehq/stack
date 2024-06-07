package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/get-momo/atlar-v1-go-client/client/transactions"
)

func (c *Client) GetV1Transactions(ctx context.Context, token string, pageSize int64) (*transactions.GetV1TransactionsOK, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "list_transactions")
	now := time.Now()
	defer f(ctx, now)

	params := transactions.GetV1TransactionsParams{
		Limit:   &pageSize,
		Context: ctx,
		Token:   &token,
	}

	return c.client.Transactions.GetV1Transactions(&params)
}

func (c *Client) GetV1TransactionsID(ctx context.Context, id string) (*transactions.GetV1TransactionsIDOK, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "list_transactions")
	now := time.Now()
	defer f(ctx, now)

	params := transactions.GetV1TransactionsIDParams{
		Context: ctx,
		ID:      id,
	}

	return c.client.Transactions.GetV1TransactionsID(&params)
}
