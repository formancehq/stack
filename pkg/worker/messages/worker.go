package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/go-libs/sharedlogging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/kafka"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/twmb/franz-go/pkg/kgo"
)

type WorkerMessages struct {
	httpClient *http.Client
	store      storage.Store

	kafkaClient kafka.Client
	kafkaTopics []string

	retriesSchedule []time.Duration

	stopChan chan chan struct{}
}

func NewWorkerMessages(store storage.Store, httpClient *http.Client, retriesSchedule []time.Duration) (*WorkerMessages, error) {
	kafkaClient, kafkaTopics, err := kafka.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "kafka.NewClient")
	}

	return &WorkerMessages{
		httpClient:      httpClient,
		store:           store,
		kafkaClient:     kafkaClient,
		kafkaTopics:     kafkaTopics,
		retriesSchedule: retriesSchedule,
		stopChan:        make(chan chan struct{}),
	}, nil
}

func (w *WorkerMessages) Run(ctx context.Context) error {
	msgChan := make(chan *kgo.Record)
	errChan := make(chan error)
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go fetchMessages(ctxWithCancel, w.kafkaClient, msgChan, errChan)

	for {
		select {
		case ch := <-w.stopChan:
			sharedlogging.GetLogger(ctx).Debug("workerMessages: received from stopChan")
			close(ch)
			return nil
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debugf("workerMessages: context done: %s", ctx.Err())
			return nil
		case err := <-errChan:
			return errors.Wrap(err, "kafka.WorkerMessages.fetchMessages")
		case msg := <-msgChan:
			ctx = sharedlogging.ContextWithLogger(ctx,
				sharedlogging.GetLogger(ctx).WithFields(map[string]any{
					"offset": msg.Offset,
				}))
			sharedlogging.GetLogger(ctx).WithFields(map[string]any{
				"time":      msg.Timestamp.UTC().Format(time.RFC3339),
				"partition": msg.Partition,
				"headers":   msg.Headers,
			}).Debug("workerMessages: new kafka message fetched")

			w.kafkaClient.PauseFetchTopics(w.kafkaTopics...)

			if err := w.processMessage(ctx, msg.Value); err != nil {
				return errors.Wrap(err, "worker.WorkerMessages.processMessage")
			}

			w.kafkaClient.ResumeFetchTopics(w.kafkaTopics...)
		}
	}
}

func (w *WorkerMessages) Stop(ctx context.Context) {
	ch := make(chan struct{})
	select {
	case <-ctx.Done():
		sharedlogging.GetLogger(ctx).Debugf("workerMessages stopped: context done: %s", ctx.Err())
		return
	case w.stopChan <- ch:
		select {
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debugf("workerMessages stopped via stopChan: context done: %s", ctx.Err())
			return
		case <-ch:
			sharedlogging.GetLogger(ctx).Debug("workerMessages stopped via stopChan")
		}
	}
}

func fetchMessages(ctx context.Context, kafkaClient kafka.Client, msgChan chan *kgo.Record, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fetches := kafkaClient.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				sharedlogging.GetLogger(ctx).Errorf("POLL: %+v", errs)
				for _, err := range errs {
					select {
					case <-ctx.Done():
						return
					default:
						errChan <- fmt.Errorf(
							"kafka.Client.PollRecords: topic: %s: partition: %d: %w",
							err.Topic, err.Partition, err.Err)
					}
				}
			}

			var records []*kgo.Record
			fetches.EachRecord(func(record *kgo.Record) {
				msgChan <- record
				records = append(records, record)
			})
		}
	}
}

func (w *WorkerMessages) processMessage(ctx context.Context, msgValue []byte) error {
	var ev webhooks.EventMessage
	if err := json.Unmarshal(msgValue, &ev); err != nil {
		return errors.Wrap(err, "json.Unmarshal event message")
	}

	eventApp := strings.ToLower(ev.App)
	eventType := strings.ToLower(ev.Type)

	if eventApp == "" {
		ev.Type = eventType
	} else {
		ev.Type = strings.Join([]string{eventApp, eventType}, ".")
	}

	filter := map[string]any{webhooks.KeyEventTypes: ev.Type}
	sharedlogging.GetLogger(ctx).Debugf("searching configs with filter: %+v", filter)
	cfgs, err := w.store.FindManyConfigs(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "storage.store.FindManyConfigs")
	}

	for _, cfg := range cfgs {
		sharedlogging.GetLogger(ctx).Debugf("found one config: %+v", cfg)
		data, err := json.Marshal(ev)
		if err != nil {
			return errors.Wrap(err, "json.Marshal event message")
		}

		attempt, err := webhooks.MakeAttempt(ctx, w.httpClient, w.retriesSchedule,
			uuid.NewString(), 0, cfg, data)
		if err != nil {
			return errors.Wrap(err, "sending webhook")
		}

		if attempt.Status == webhooks.StatusAttemptSuccess {
			sharedlogging.GetLogger(ctx).Infof(
				"webhook sent with ID %s to %s of type %s",
				attempt.WebhookID, cfg.Endpoint, ev.Type)
		}

		if err := w.store.InsertOneAttempt(ctx, attempt); err != nil {
			return errors.Wrap(err, "storage.store.InsertOneAttempt")
		}
	}

	return nil
}
