package client

import (
	"context"

	"github.com/stripe/stripe-go/v79"
)

func (c *client) GetExternalAccounts(
	ctx context.Context,
	accountID *string,
	lastID *string,
	pageSize int64,
) ([]*stripe.BankAccount, bool, error) {
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

	itr := c.bankAccountClient.List(&stripe.BankAccountListParams{
		Account:    accountID,
		ListParams: filters,
	})
	if err := itr.Err(); err != nil {
		return nil, false, handleError(itr.Iter)
	}
	return itr.BankAccountList().Data, itr.BankAccountList().ListMeta.HasMore, nil
}
