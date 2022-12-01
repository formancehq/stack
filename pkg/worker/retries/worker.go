package retries

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/go-libs/sharedlogging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/pkg/errors"
)

type WorkerRetries struct {
	httpClient *http.Client
	store      storage.Store

	retriesCron     time.Duration
	retriesSchedule []time.Duration

	stopChan chan chan struct{}
}

func NewWorkerRetries(store storage.Store, httpClient *http.Client, retriesCron time.Duration, retriesSchedule []time.Duration) (*WorkerRetries, error) {
	return &WorkerRetries{
		httpClient:      httpClient,
		store:           store,
		retriesCron:     retriesCron,
		retriesSchedule: retriesSchedule,
		stopChan:        make(chan chan struct{}),
	}, nil
}

func (w *WorkerRetries) Run(ctx context.Context) error {
	errChan := make(chan error)
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go w.attemptRetries(ctxWithCancel, errChan)

	for {
		select {
		case ch := <-w.stopChan:
			sharedlogging.GetLogger(ctx).Debug("workerRetries: received from stopChan")
			close(ch)
			return nil
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debugf("workerRetries: context done: %s", ctx.Err())
			return nil
		case err := <-errChan:
			return errors.Wrap(err, "kafka.WorkerRetries.attemptRetries")
		}
	}
}

func (w *WorkerRetries) Stop(ctx context.Context) {
	ch := make(chan struct{})
	select {
	case <-ctx.Done():
		sharedlogging.GetLogger(ctx).Debugf("workerRetries stopped: context done: %s", ctx.Err())
		return
	case w.stopChan <- ch:
		select {
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debugf("workerRetries stopped via stopChan: context done: %s", ctx.Err())
			return
		case <-ch:
			sharedlogging.GetLogger(ctx).Debug("workerRetries stopped via stopChan")
		}
	}
}

var ErrNoAttemptsFound = errors.New("attemptRetries: no attempts found")

func (w *WorkerRetries) attemptRetries(ctx context.Context, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Find all webhookIDs ready to be retried
			filter := map[string]any{
				webhooks.KeyStatus:         webhooks.StatusAttemptToRetry,
				webhooks.KeyNextRetryAfter: map[string]any{"$lt": time.Now().UTC()},
			}
			ids, err := w.store.FindDistinctWebhookIDs(ctx, filter)
			if err != nil {
				errChan <- errors.Wrap(err, "storage.Store.FindDistinctWebhookIDs to retry")
				continue
			} else {
				sharedlogging.GetLogger(ctx).Debugf("found %d distinct webhookIDs to retry: %+v", len(ids), ids)
			}

			for _, id := range ids {
				filter[webhooks.KeyWebhookID] = id
				atts, err := w.store.FindManyAttempts(ctx, filter)
				if err != nil {
					errChan <- errors.Wrap(err, "storage.Store.FindManyAttempts")
					continue
				}
				if len(atts) == 0 {
					errChan <- fmt.Errorf("%w for webhookID: %s", ErrNoAttemptsFound, id)
					continue
				}

				newAttemptNb := atts[0].RetryAttempt + 1
				attempt, err := webhooks.MakeAttempt(ctx, w.httpClient, w.retriesSchedule,
					id, newAttemptNb, atts[0].Config, []byte(atts[0].Payload))
				if err != nil {
					errChan <- errors.Wrap(err, "webhooks.MakeAttempt")
					continue
				}

				if err := w.store.InsertOneAttempt(ctx, attempt); err != nil {
					errChan <- errors.Wrap(err, "storage.Store.InsertOneAttempt retried")
					continue
				}

				if _, err := w.store.UpdateManyAttemptsStatus(ctx, id, attempt.Status); err != nil {
					if errors.Is(err, storage.ErrAttemptNotModified) {
						continue
					}
					errChan <- errors.Wrap(err, "storage.Store.UpdateManyAttemptsStatus")
					continue
				}
			}
		}

		time.Sleep(w.retriesCron)
	}
}
