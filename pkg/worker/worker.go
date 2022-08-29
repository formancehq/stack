package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/pkg/engine"
	"github.com/numary/webhooks/pkg/engine/svix"
	"github.com/numary/webhooks/pkg/storage"
	kafkago "github.com/segmentio/kafka-go"
)

type Worker struct {
	Reader Reader
	Store  storage.Store

	engine engine.Engine

	stopChan chan chan struct{}
}

func NewWorker(store storage.Store, engine svix.Engine) (*Worker, error) {
	cfg, err := NewKafkaReaderConfig()
	if err != nil {
		return nil, fmt.Errorf("NewKafkaReaderConfig: %w", err)
	}

	return &Worker{
		Reader:   kafkago.NewReader(cfg),
		Store:    store,
		engine:   engine,
		stopChan: make(chan chan struct{}),
	}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	msgChan := make(chan kafkago.Message)
	errChan := make(chan error)
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go fetchMessages(ctxWithCancel, w.Reader, msgChan, errChan)

	for {
		select {
		case ch := <-w.stopChan:
			sharedlogging.GetLogger(ctx).Debug("worker: received from stopChan")
			close(ch)
			return nil
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debugf("worker: context done: %s", ctx.Err())
			return nil
		case err := <-errChan:
			return fmt.Errorf("kafka.Worker.fetchMessages: %w", err)
		case msg := <-msgChan:
			ctx = sharedlogging.ContextWithLogger(ctx,
				sharedlogging.GetLogger(ctx).WithFields(map[string]any{
					"offset": msg.Offset,
				}))
			sharedlogging.GetLogger(ctx).WithFields(map[string]any{
				"time":      msg.Time.UTC().Format(time.RFC3339),
				"partition": msg.Partition,
				"headers":   msg.Headers,
			}).Debug("worker: new kafka message fetched")

			eventType, err := FilterMessage(msg.Value)
			if err != nil {
				return err
			}

			if err := w.engine.ProcessKafkaMessage(ctx, eventType, msg.Value); err != nil {
				return fmt.Errorf("engine.ProcessKafkaMessage: %w", err)
			}

			if err := w.Reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("kafka.Reader.CommitMessages: %w", err)
			}
		}
	}
}

func fetchMessages(ctx context.Context, reader Reader, msgChan chan kafkago.Message, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if !errors.Is(err, io.EOF) && ctx.Err() == nil {
					select {
					case errChan <- fmt.Errorf("kafka.Reader.FetchMessage: %w", err):
					case <-ctx.Done():
						return
					}
				}
				continue
			}

			select {
			case msgChan <- msg:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (w *Worker) Stop(ctx context.Context) {
	ch := make(chan struct{})
	select {
	case <-ctx.Done():
		sharedlogging.GetLogger(ctx).Debugf("worker stopped: context done: %s", ctx.Err())
		return
	case w.stopChan <- ch:
		select {
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debugf("worker stopped via stopChan: context done: %s", ctx.Err())
			return
		case <-ch:
			sharedlogging.GetLogger(ctx).Debug("worker stopped via stopChan")
		}
	}
}
