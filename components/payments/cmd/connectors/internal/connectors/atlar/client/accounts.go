package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/get-momo/atlar-v1-go-client/client/accounts"
)

func (c *Client) GetV1AccountsID(ctx context.Context, id string) (*accounts.GetV1AccountsIDOK, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "list_accounts")
	now := time.Now()
	defer f(ctx, now)

	accountsParams := accounts.GetV1AccountsIDParams{
		Context: ctx,
		ID:      id,
	}

	return c.client.Accounts.GetV1AccountsID(&accountsParams)
}

func (c *Client) GetV1Accounts(ctx context.Context, token string, pageSize int64) (*accounts.GetV1AccountsOK, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "list_accounts")
	now := time.Now()
	defer f(ctx, now)

	accountsParams := accounts.GetV1AccountsParams{
		Limit:   &pageSize,
		Context: ctx,
		Token:   &token,
	}

	return c.client.Accounts.GetV1Accounts(&accountsParams)
}
