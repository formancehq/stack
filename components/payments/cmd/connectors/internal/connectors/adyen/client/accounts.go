package client

import (
	"context"
	"fmt"
	"time"

	"github.com/adyen/adyen-go-api-library/v7/src/management"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

func (c *Client) GetMerchantAccounts(ctx context.Context, pageNumber, pageSize int32) ([]management.Merchant, error) {
	f := connectors.ClientMetrics(ctx, "adyen", "list_merchant_accounts")
	now := time.Now()
	defer f(ctx, now)

	listMerchantsResponse, raw, err := c.client.Management().AccountMerchantLevelApi.ListMerchantAccounts(
		ctx,
		c.client.Management().AccountMerchantLevelApi.ListMerchantAccountsInput().PageNumber(pageNumber).PageSize(pageSize),
	)
	if err != nil {
		return nil, err
	}

	if raw.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get merchant accounts: %d", raw.StatusCode)
	}

	return listMerchantsResponse.Data, nil
}
