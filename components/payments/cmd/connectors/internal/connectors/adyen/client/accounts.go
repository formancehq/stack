package client

import (
	"context"
	"fmt"

	"github.com/adyen/adyen-go-api-library/v7/src/management"
)

func (c *Client) GetMerchantAccounts(ctx context.Context, pageNumber, pageSize int32) ([]management.Merchant, error) {
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
