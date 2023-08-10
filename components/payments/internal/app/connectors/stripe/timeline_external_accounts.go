package stripe

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/stripe/stripe-go/v72"
)

//nolint:tagliatelle // allow different styled tags in client
type ExternalAccountsListResponse struct {
	HasMore bool                      `json:"has_more"`
	Data    []*stripe.ExternalAccount `json:"data"`
}

func (tl *Timeline) doExternalAccountsRequest(ctx context.Context, queryParams url.Values,
	to *[]*stripe.ExternalAccount,
) (bool, error) {
	options := make([]ClientOption, 0)
	options = append(options, QueryParam("limit", fmt.Sprintf("%d", tl.config.PageSize)))

	for k, v := range queryParams {
		options = append(options, QueryParam(k, v[0]))
	}

	txs, hasMore, err := tl.client.ExternalAccounts(ctx, options...)
	if err != nil {
		return false, err
	}

	*to = txs

	return hasMore, nil
}

func (tl *Timeline) initExternalAccounts(ctx context.Context) error {
	ret := make([]*stripe.ExternalAccount, 0)

	_, err := tl.doExternalAccountsRequest(ctx, url.Values{}, &ret)
	if err != nil {
		return err
	}

	if len(ret) > 0 {
		tl.firstIDAfterStartingAt = ret[0].ID
	}

	return nil
}

func (tl *Timeline) ExternalAccountsTail(ctx context.Context, to *[]*stripe.ExternalAccount) (bool, TimelineState, func(), error) {
	queryParams := url.Values{}

	switch {
	case tl.state.OldestID != "":
		queryParams.Set("starting_after", tl.state.OldestID)
	default:
	}

	hasMore, err := tl.doExternalAccountsRequest(ctx, queryParams, to)
	if err != nil {
		return false, TimelineState{}, nil, err
	}

	futureState := tl.state

	if len(*to) > 0 {
		lastItem := (*to)[len(*to)-1]
		futureState.OldestID = lastItem.ID
		oldestDate := time.Unix(lastItem.BankAccount.Account.Created, 0)
		futureState.OldestDate = &oldestDate

		if futureState.MoreRecentID == "" {
			firstItem := (*to)[0]
			futureState.MoreRecentID = firstItem.ID
			moreRecentDate := time.Unix(firstItem.BankAccount.Account.Created, 0)
			futureState.MoreRecentDate = &moreRecentDate
		}
	}

	futureState.NoMoreHistory = !hasMore

	return hasMore, futureState, func() {
		tl.state = futureState
	}, nil
}

func (tl *Timeline) ExternalAccountsHead(ctx context.Context, to *[]*stripe.ExternalAccount) (bool, TimelineState, func(), error) {
	if tl.firstIDAfterStartingAt == "" && tl.state.MoreRecentID == "" {
		err := tl.initExternalAccounts(ctx)
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

	hasMore, err := tl.doExternalAccountsRequest(ctx, queryParams, to)
	if err != nil {
		return false, TimelineState{}, nil, err
	}

	futureState := tl.state

	if len(*to) > 0 {
		firstItem := (*to)[0]
		futureState.MoreRecentID = firstItem.ID
		moreRecentDate := time.Unix(firstItem.BankAccount.Account.Created, 0)
		futureState.MoreRecentDate = &moreRecentDate

		if futureState.OldestID == "" {
			lastItem := (*to)[len(*to)-1]
			futureState.OldestID = lastItem.ID
			oldestDate := time.Unix(lastItem.BankAccount.Account.Created, 0)
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
