package stripe

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/stripe/client"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v72"
)

func TestTimeline(t *testing.T) {
	t.Parallel()

	mock := NewClientMock(t, true)
	ref := time.Now()
	timeline := NewTimeline(mock, TimelineConfig{
		PageSize: 2,
	}, TimelineState{}, WithStartingAt(ref))

	tx1 := &stripe.BalanceTransaction{
		ID:      "tx1",
		Created: ref.Add(-time.Minute).Unix(),
	}

	tx2 := &stripe.BalanceTransaction{
		ID:      "tx2",
		Created: ref.Add(-2 * time.Minute).Unix(),
	}

	mock.Expect().
		Limit(2).
		CreatedLte(ref).
		RespondsWith(true, tx1, tx2)

	ret := make([]*stripe.BalanceTransaction, 0)
	hasMore, state, commit, err := timeline.TransactionsTail(context.TODO(), &ret)
	require.NoError(t, err)
	require.True(t, hasMore)
	require.Equal(t, TimelineState{
		OldestID:       "tx2",
		OldestDate:     client.DatePtr(time.Unix(tx2.Created, 0)),
		MoreRecentID:   "tx1",
		MoreRecentDate: client.DatePtr(time.Unix(tx1.Created, 0)),
		NoMoreHistory:  false,
	}, state)

	commit()

	tx3 := &stripe.BalanceTransaction{
		ID:      "tx3",
		Created: ref.Add(-3 * time.Minute).Unix(),
	}

	mock.Expect().Limit(2).StartingAfter(tx2.ID).RespondsWith(false, tx3)

	hasMore, state, _, err = timeline.TransactionsTail(context.TODO(), &ret)
	require.NoError(t, err)
	require.False(t, hasMore)
	require.Equal(t, TimelineState{
		OldestID:       "tx3",
		OldestDate:     client.DatePtr(time.Unix(tx3.Created, 0)),
		MoreRecentID:   "tx1",
		MoreRecentDate: client.DatePtr(time.Unix(tx1.Created, 0)),
		NoMoreHistory:  true,
	}, state)
}
