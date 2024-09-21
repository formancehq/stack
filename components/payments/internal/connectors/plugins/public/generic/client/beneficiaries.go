package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/genericclient"
)

func (c *Client) ListBeneficiaries(ctx context.Context, page, pageSize int64, createdAtFrom time.Time) ([]genericclient.Beneficiary, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "generic", "list_beneficiaries")
	// now := time.Now()
	// defer f(ctx, now)

	req := c.apiClient.DefaultApi.
		GetBeneficiaries(ctx).
		Page(page).
		PageSize(pageSize)

	if !createdAtFrom.IsZero() {
		req = req.CreatedAtFrom(createdAtFrom)
	}

	beneficiaries, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	return beneficiaries, nil
}
