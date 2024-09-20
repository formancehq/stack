package client

import (
	"context"

	"github.com/stripe/stripe-go/v79"
)

func (c *client) GetAccounts(ctx context.Context, lastID *string, pageSize int64) ([]*stripe.Account, bool, error) {
	// TODO
	//	f := connectors.ClientMetrics(ctx, "stripe", "list_accounts")
	//	now := time.Now()
	//	defer f(ctx, now)
	filters := stripe.ListParams{
		Limit: &pageSize,
	}
	if lastID == nil {
		filters.StartingAfter = lastID
	}

	itr := c.accountClient.List(&stripe.AccountListParams{ListParams: filters})
	if err := itr.Err(); err != nil {
		return nil, false, handleError(itr.Iter)
	}
	return itr.AccountList().Data, itr.AccountList().ListMeta.HasMore, nil
}
