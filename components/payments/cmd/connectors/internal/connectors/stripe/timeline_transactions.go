package stripe

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/stripe/stripe-go/v72"
)

func (tl *Timeline) doTransactionsRequest(ctx context.Context, queryParams url.Values,
	to *[]*stripe.BalanceTransaction,
) (bool, error) {
	options := make([]client.ClientOption, 0)
	options = append(options, client.QueryParam("limit", fmt.Sprintf("%d", tl.config.PageSize)))
	options = append(options, client.QueryParam("expand[]", "data.source"))

	for k, v := range queryParams {
		options = append(options, client.QueryParam(k, v[0]))
	}

	txs, hasMore, err := tl.client.BalanceTransactions(ctx, options...)
	if err != nil {
		return false, err
	}

	*to = txs

	return hasMore, nil
}

func (tl *Timeline) initTransactions(ctx context.Context) error {
	ret := make([]*stripe.BalanceTransaction, 0)
	params := url.Values{}
	params.Set("limit", "1")
	params.Set("created[lt]", fmt.Sprintf("%d", tl.startingAt.Unix()))

	_, err := tl.doTransactionsRequest(ctx, params, &ret)
	if err != nil {
		return err
	}

	if len(ret) > 0 {
		tl.firstIDAfterStartingAt = ret[0].ID
	}

	return nil
}

func (tl *Timeline) TransactionsTail(ctx context.Context, to *[]*stripe.BalanceTransaction) (bool, TimelineState, func(), error) {
	queryParams := url.Values{}

	switch {
	case tl.state.OldestID != "":
		queryParams.Set("starting_after", tl.state.OldestID)
	default:
		queryParams.Set("created[lte]", fmt.Sprintf("%d", tl.startingAt.Unix()))
	}

	hasMore, err := tl.doTransactionsRequest(ctx, queryParams, to)
	if err != nil {
		return false, TimelineState{}, nil, err
	}

	futureState := tl.state

	if len(*to) > 0 {
		lastItem := (*to)[len(*to)-1]
		futureState.OldestID = lastItem.ID
		oldestDate := time.Unix(lastItem.Created, 0)
		futureState.OldestDate = &oldestDate

		if futureState.MoreRecentID == "" {
			firstItem := (*to)[0]
			futureState.MoreRecentID = firstItem.ID
			moreRecentDate := time.Unix(firstItem.Created, 0)
			futureState.MoreRecentDate = &moreRecentDate
		}
	}

	futureState.NoMoreHistory = !hasMore

	return hasMore, futureState, func() {
		tl.state = futureState
	}, nil
}

func (tl *Timeline) TransactionsHead(ctx context.Context, to *[]*stripe.BalanceTransaction) (bool, TimelineState, func(), error) {
	if tl.firstIDAfterStartingAt == "" && tl.state.MoreRecentID == "" {
		err := tl.initTransactions(ctx)
		if err != nil {
			return false, TimelineState{}, nil, err
		}

		if tl.firstIDAfterStartingAt == "" {
			return false, TimelineState{
				NoMoreHistory: true,
			}, func() {}, nil
		}
	}

	queryParams := url.Values{}

	switch {
	case tl.state.MoreRecentID != "":
		queryParams.Set("ending_before", tl.state.MoreRecentID)
	case tl.firstIDAfterStartingAt != "":
		queryParams.Set("ending_before", tl.firstIDAfterStartingAt)
	}

	hasMore, err := tl.doTransactionsRequest(ctx, queryParams, to)
	if err != nil {
		return false, TimelineState{}, nil, err
	}

	futureState := tl.state

	if len(*to) > 0 {
		firstItem := (*to)[0]
		futureState.MoreRecentID = firstItem.ID
		moreRecentDate := time.Unix(firstItem.Created, 0)
		futureState.MoreRecentDate = &moreRecentDate

		if futureState.OldestID == "" {
			lastItem := (*to)[len(*to)-1]
			futureState.OldestID = lastItem.ID
			oldestDate := time.Unix(lastItem.Created, 0)
			futureState.OldestDate = &oldestDate
		}
	}

	futureState.NoMoreHistory = !hasMore

	for i, j := 0, len(*to)-1; i < j; i, j = i+1, j-1 {
		(*to)[i], (*to)[j] = (*to)[j], (*to)[i]
	}

	return hasMore, futureState, func() {
		tl.state = futureState
	}, nil
}
