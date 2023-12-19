package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/get-momo/atlar-v1-go-client/client/external_accounts"
)

func (c *Client) GetV1ExternalAccountsID(ctx context.Context, externalAccountID string) (*external_accounts.GetV1ExternalAccountsIDOK, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "get_external_account")
	now := time.Now()
	defer f(ctx, now)

	getExternalAccountParams := external_accounts.GetV1ExternalAccountsIDParams{
		Context: ctx,
		ID:      externalAccountID,
	}

	externalAccountResponse, err := c.client.ExternalAccounts.GetV1ExternalAccountsID(&getExternalAccountParams)
	if err != nil {
		return nil, err
	}

	return externalAccountResponse, nil
}

func (c *Client) GetV1ExternalAccounts(ctx context.Context, token string, pageSize int64) (*external_accounts.GetV1ExternalAccountsOK, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "list_external_accounts")
	now := time.Now()
	defer f(ctx, now)

	externalAccountsParams := external_accounts.GetV1ExternalAccountsParams{
		Limit:   &pageSize,
		Context: ctx,
		Token:   &token,
	}

	return c.client.ExternalAccounts.GetV1ExternalAccounts(&externalAccountsParams)
}
