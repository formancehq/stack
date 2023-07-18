package stripe

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
	"golang.org/x/sync/semaphore"
)

type TimelineTriggerType string

const (
	TimelineTriggerTypeTransactions TimelineTriggerType = "transactions"
	TimelineTriggerTypeAccounts     TimelineTriggerType = "accounts"
)

func NewTimelineTrigger(
	logger logging.Logger,
	ingester Ingester,
	timeline *Timeline,
	timelineType TimelineTriggerType,
) *TimelineTrigger {
	return &TimelineTrigger{
		logger: logger.WithFields(map[string]interface{}{
			"component": "timeline-trigger",
		}),
		ingester:     ingester,
		timeline:     timeline,
		timelineType: timelineType,
		sem:          semaphore.NewWeighted(1),
	}
}

type TimelineTrigger struct {
	logger       logging.Logger
	ingester     Ingester
	timeline     *Timeline
	timelineType TimelineTriggerType
	sem          *semaphore.Weighted
	cancel       func()
}

func (t *TimelineTrigger) Fetch(ctx context.Context) error {
	if t.sem.TryAcquire(1) {
		defer t.sem.Release(1)

		ctx, t.cancel = context.WithCancel(ctx)
		if !t.timeline.State().NoMoreHistory {
			if err := t.fetch(ctx, true); err != nil {
				return err
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := t.fetch(ctx, false); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TimelineTrigger) Cancel(ctx context.Context) {
	if t.cancel != nil {
		t.cancel()

		err := t.sem.Acquire(ctx, 1)
		if err != nil {
			panic(err)
		}

		t.sem.Release(1)
	}
}

func (t *TimelineTrigger) fetch(ctx context.Context, tail bool) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			hasMore, err := t.triggerPage(ctx, tail)
			if err != nil {
				return errors.Wrap(err, "error triggering tail page")
			}

			if !hasMore {
				return nil
			}
		}
	}
}

func (t *TimelineTrigger) triggerPage(ctx context.Context, tail bool) (bool, error) {
	logger := t.logger.WithFields(map[string]interface{}{
		"tail": tail,
	})

	logger.Debugf("Trigger page")

	var hasMore bool
	switch t.timelineType {
	case TimelineTriggerTypeTransactions:
		ret := make([]*stripe.BalanceTransaction, 0)
		method := t.timeline.TransactionsHead
		if tail {
			method = t.timeline.TransactionsTail
		}

		more, futureState, commitFn, err := method(ctx, &ret)
		if err != nil {
			return false, errors.Wrap(err, "fetching timeline")
		}
		hasMore = more

		logger.Debug("Ingest transactions batch")

		if len(ret) > 0 {
			err = t.ingester.IngestTransactions(ctx, ret, futureState, tail)
			if err != nil {
				return false, errors.Wrap(err, "ingesting batch")
			}
		}

		commitFn()

	case TimelineTriggerTypeAccounts:
		ret := make([]*stripe.Account, 0)
		method := t.timeline.AccountsHead
		if tail {
			method = t.timeline.AccountsTail
		}

		more, futureState, commitFn, err := method(ctx, &ret)
		if err != nil {
			return false, errors.Wrap(err, "fetching timeline")
		}
		hasMore = more

		logger.Debug("Ingest transactions batch")

		if len(ret) > 0 {
			err = t.ingester.IngestAccounts(ctx, ret, futureState, tail)
			if err != nil {
				return false, errors.Wrap(err, "ingesting batch")
			}
		}

		commitFn()
	}

	return hasMore, nil
}
