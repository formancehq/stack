package client

import (
	"fmt"

	"github.com/stripe/stripe-go/v79"
)

// Timeline allows the client to navigate the backlog and decide whether to fetch
// historical or recently added data
type Timeline struct {
	LatestID      string `json:"latest_id"`
	BacklogCursor string `json:"backlog_cursor"`
}

func (t Timeline) IsCaughtUp() bool {
	return t.LatestID != ""
}

func scanForOldest(
	timeline Timeline,
	pageSize int64,
	listFn func(stripe.ListParams) (stripe.ListContainer, error),
) (interface{}, Timeline, bool, error) {
	filters := stripe.ListParams{
		Limit:  limit(pageSize, 0),
		Single: true, // turn off autopagination
	}
	if timeline.BacklogCursor != "" {
		filters.StartingAfter = &timeline.BacklogCursor
	}

	var oldest interface{}
	var oldestID string

	list, err := listFn(filters)
	if err != nil {
		return oldest, timeline, false, err
	}
	hasMore := list.GetListMeta().HasMore

	switch v := list.(type) {
	case *stripe.AccountList:
		if len(v.Data) == 0 {
			return oldest, timeline, hasMore, nil
		}
		account := v.Data[len(v.Data)-1]
		oldest = account
		oldestID = account.ID

	case *stripe.BankAccountList:
		if len(v.Data) == 0 {
			return oldest, timeline, hasMore, nil
		}
		account := v.Data[len(v.Data)-1]
		oldest = account
		oldestID = account.ID

	case *stripe.BalanceTransactionList:
		if len(v.Data) == 0 {
			return oldest, timeline, hasMore, nil
		}
		trx := v.Data[len(v.Data)-1]
		oldest = trx
		oldestID = trx.ID
	default:
		return nil, timeline, hasMore, fmt.Errorf("failed to fetch backlog for type %T", list)
	}

	// we haven't found the oldest yet
	if hasMore {
		timeline.BacklogCursor = oldestID
		return nil, timeline, hasMore, nil
	}
	timeline.LatestID = oldestID
	return oldest, timeline, hasMore, nil
}
