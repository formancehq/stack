package worker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Retrier struct {
	httpClient *http.Client
	store      storage.Store

	retriesCron     time.Duration
	retriesSchedule []time.Duration

	stopChan chan chan struct{}
}

func NewRetrier(store storage.Store, httpClient *http.Client, retriesCron time.Duration, retriesSchedule []time.Duration) (*Retrier, error) {
	return &Retrier{
		httpClient:      httpClient,
		store:           store,
		retriesCron:     retriesCron,
		retriesSchedule: retriesSchedule,
		stopChan:        make(chan chan struct{}),
	}, nil
}

func (w *Retrier) Run(ctx context.Context) error {
	errChan := make(chan error)
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go w.attemptRetries(ctxWithCancel, errChan)

	for {
		select {
		case ch := <-w.stopChan:
			logging.FromContext(ctx).Debug("worker: received from stopChan")
			close(ch)
			return nil
		case <-ctx.Done():
			logging.FromContext(ctx).Debugf("worker: context done: %s", ctx.Err())
			return nil
		case err := <-errChan:
			return errors.Wrap(err, "kafka.Retrier")
		}
	}
}

func (w *Retrier) Stop(ctx context.Context) {
	ch := make(chan struct{})
	select {
	case <-ctx.Done():
		logging.FromContext(ctx).Debugf("worker stopped: context done: %s", ctx.Err())
		return
	case w.stopChan <- ch:
		select {
		case <-ctx.Done():
			logging.FromContext(ctx).Debugf("worker stopped via stopChan: context done: %s", ctx.Err())
			return
		case <-ch:
			logging.FromContext(ctx).Debug("worker stopped via stopChan")
		}
	default:
		logging.FromContext(ctx).Debug("trying to stop worker: no communication")
	}
}

var ErrNoAttemptsFound = errors.New("attemptRetries: no attempts found")

func (w *Retrier) attemptRetries(ctx context.Context, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Find all webhookIDs ready to be retried
			webhookIDs, err := w.store.FindWebhookIDsToRetry(ctx)
			if err != nil {
				errChan <- errors.Wrap(err, "storage.Store.FindWebhookIDsToRetry")
				continue
			} else {
				logging.FromContext(ctx).Debugf(
					"found %d distinct webhookIDs to retry: %+v", len(webhookIDs), webhookIDs)
			}

			for _, webhookID := range webhookIDs {
				atts, err := w.store.FindAttemptsToRetryByWebhookID(ctx, webhookID)
				if err != nil {
					errChan <- errors.Wrap(err, "storage.Store.FindAttemptsToRetryByWebhookID")
					continue
				}
				if len(atts) == 0 {
					errChan <- fmt.Errorf("%w for webhookID: %s", ErrNoAttemptsFound, webhookID)
					continue
				}

				newAttemptNb := atts[0].RetryAttempt + 1
				attempt, err := webhooks.MakeAttempt(ctx, w.httpClient, w.retriesSchedule, uuid.NewString(),
					webhookID, newAttemptNb, atts[0].Config, []byte(atts[0].Payload), false)
				if err != nil {
					errChan <- errors.Wrap(err, "webhooks.MakeAttempt")
					continue
				}

				if err := w.store.InsertOneAttempt(ctx, attempt); err != nil {
					errChan <- errors.Wrap(err, "storage.Store.InsertOneAttempt retried")
					continue
				}

				if _, err := w.store.UpdateAttemptsStatus(ctx, webhookID, attempt.Status); err != nil {
					if errors.Is(err, storage.ErrAttemptsNotModified) {
						continue
					}
					errChan <- errors.Wrap(err, "storage.Store.UpdateAttemptsStatus")
					continue
				}
			}
		}

		time.Sleep(w.retriesCron)
	}
}
