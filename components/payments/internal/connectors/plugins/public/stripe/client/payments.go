package client

import (
	"context"

	"github.com/stripe/stripe-go/v79"
)

func (c *client) GetPayments(ctx context.Context, accountID *string, lastID *string, pageSize int64) ([]*stripe.BalanceTransaction, bool, error) {
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

	itr := c.balanceTransactionClient.List(&stripe.BalanceTransactionListParams{ListParams: filters})
	if err := itr.Err(); err != nil {
		return nil, false, handleError(itr.Iter)
	}
	return itr.BalanceTransactionList().Data, itr.BalanceTransactionList().ListMeta.HasMore, nil
}
