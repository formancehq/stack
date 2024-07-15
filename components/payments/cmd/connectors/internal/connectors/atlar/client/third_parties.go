package client

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/get-momo/atlar-v1-go-client/client/third_parties"
)

func (c *Client) GetV1BetaThirdPartiesID(ctx context.Context, id string) (*third_parties.GetV1betaThirdPartiesIDOK, error) {
	f := connectors.ClientMetrics(ctx, "atlar", "list_third_parties")
	now := time.Now()
	defer f(ctx, now)

	params := third_parties.GetV1betaThirdPartiesIDParams{
		Context: ctx,
		ID:      id,
	}

	return c.client.ThirdParties.GetV1betaThirdPartiesID(&params)
}
