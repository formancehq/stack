package client

import (
	"context"

	"github.com/stripe/stripe-go/v79"
)

func (c *client) GetExternalAccounts(
	ctx context.Context,
	accountID string,
	timeline Timeline,
	pageSize int64,
) (results []*stripe.BankAccount, _ Timeline, hasMore bool, err error) {
	// TODO
	//	f := connectors.ClientMetrics(ctx, "stripe", "list_accounts")
	//	now := time.Now()
	//	defer f(ctx, now)

	results = make([]*stripe.BankAccount, 0, int(pageSize))

	// return 0 results because this endpoint cannot be used for root account
	if accountID == "" {
		return results, timeline, false, nil
	}

	if !timeline.IsCaughtUp() {
		var oldest interface{}
		oldest, timeline, hasMore, err = scanForOldest(timeline, pageSize, func(params stripe.ListParams) (stripe.ListContainer, error) {
			itr := c.bankAccountClient.List(&stripe.BankAccountListParams{
				Account:    &accountID,
				ListParams: params,
			})
			return itr.BankAccountList(), itr.Err()
		})
		if err != nil {
			return results, timeline, false, err
		}
		// either there are no records or we haven't found the start yet
		if !timeline.IsCaughtUp() {
			return results, timeline, hasMore, nil
		}
		results = append(results, oldest.(*stripe.BankAccount))
	}

	itr := c.bankAccountClient.List(&stripe.BankAccountListParams{
		Account: &accountID,
		ListParams: stripe.ListParams{
			Limit:        &pageSize,
			EndingBefore: &timeline.LatestID,
		},
	})
	if err := itr.Err(); err != nil {
		return nil, timeline, false, handleError(itr.Iter)
	}
	results = append(results, itr.BankAccountList().Data...)
	timeline.LatestID = results[len(results)-1].ID
	return results, timeline, itr.BankAccountList().ListMeta.HasMore, nil
}
