package client

import (
	"context"

	"github.com/stripe/stripe-go/v79"
)

func (c *client) GetAccounts(
	ctx context.Context,
	timeline Timeline,
	pageSize int64,
) (results []*stripe.Account, _ Timeline, hasMore bool, err error) {
	// TODO
	//	f := connectors.ClientMetrics(ctx, "stripe", "list_accounts")
	//	now := time.Now()
	//	defer f(ctx, now)

	results = make([]*stripe.Account, 0, int(pageSize))

	if !timeline.IsCaughtUp() {
		var oldest interface{}
		oldest, timeline, hasMore, err = scanForOldest(timeline, pageSize, func(params stripe.ListParams) (stripe.ListContainer, error) {
			itr := c.accountClient.List(&stripe.AccountListParams{ListParams: params})
			return itr.AccountList(), itr.Err()
		})
		if err != nil {
			return results, timeline, false, err
		}
		// either there are no records or we haven't found the start yet
		if !timeline.IsCaughtUp() {
			return results, timeline, hasMore, nil
		}
		results = append(results, oldest.(*stripe.Account))
	}

	filters := stripe.ListParams{
		Limit:        limit(pageSize, len(results)),
		EndingBefore: &timeline.LatestID,
		Single:       true, // turn off autopagination
	}

	itr := c.accountClient.List(&stripe.AccountListParams{ListParams: filters})
	results = append(results, itr.AccountList().Data...)
	timeline.LatestID = results[len(results)-1].ID
	return results, timeline, itr.AccountList().ListMeta.HasMore, itr.Err()
}
