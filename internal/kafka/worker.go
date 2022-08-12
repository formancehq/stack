package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/internal/storage"
	svixgo "github.com/svix/svix-webhooks/go"
)

type Event struct {
	Date    time.Time      `json:"date"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

type Worker struct {
	reader     Reader
	store      storage.Store
	svixClient *svixgo.Svix
	svixAppId  string

	eventChan   chan *Event
	messageChan chan *svixgo.MessageOut
	errChan     chan error
	stopChan    chan chan struct{}
}

func NewWorker(reader Reader, store storage.Store, svixClient *svixgo.Svix, svixAppId string) *Worker {
	return &Worker{
		reader:      reader,
		store:       store,
		svixClient:  svixClient,
		svixAppId:   svixAppId,
		eventChan:   make(chan *Event),
		messageChan: make(chan *svixgo.MessageOut),
		errChan:     make(chan error),
		stopChan:    make(chan chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) {
	sharedlogging.GetLogger(ctx).Debug("worker started")

	for {
		select {
		case ch := <-w.stopChan:
			sharedlogging.GetLogger(ctx).Debug("stopping worker...")
			close(ch)
			w.closeChannels()
			return
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debug("worker: context done")
			w.errChan <- ctx.Err()
			w.closeChannels()
			return
		default:
			sharedlogging.GetLogger(ctx).Debug("worker: fetching message...")
		}

		m, err := w.reader.FetchMessage(ctx)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				w.errChan <- err
				sharedlogging.GetLogger(ctx).Errorf("kafka.Reader.FetchMessage: %s", err)
			}
			continue
		}

		ctx := sharedlogging.ContextWithLogger(ctx,
			sharedlogging.GetLogger(ctx).WithFields(map[string]any{
				"offset": m.Offset,
			}))
		sharedlogging.GetLogger(ctx).WithFields(map[string]any{
			"time":      m.Time.UTC().Format(time.RFC3339),
			"partition": m.Partition,
			"data":      string(m.Value),
			"headers":   m.Headers,
		}).Debug("new kafka message fetched")

		ev := Event{}
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			w.errChan <- err
			sharedlogging.GetLogger(ctx).Errorf("json.Unmarshal: %s", err)
			continue
		}

		sharedlogging.GetLogger(ctx).Debug("event sent in the channel")
		w.eventChan <- &ev

		if out, err := w.svixClient.EventType.Get(ev.Type); err == nil {
			dumpOut := spew.Sdump(out)
			sharedlogging.GetLogger(ctx).Debug("svix.Svix.EventType.Get: ",
				"ev.Type: ", ev.Type, "dumpOut: ", dumpOut)
			id := uuid.New().String()
			messageIn := &svixgo.MessageIn{
				EventType: ev.Type,
				EventId:   *svixgo.NullableString(id),
				Payload:   ev.Payload,
			}
			options := &svixgo.PostOptions{IdempotencyKey: &id}
			dumpIn := spew.Sdump(
				"svix appId: ", w.svixAppId,
				"svix.MessageIn: ", messageIn,
				"svix.PostOptions: ", options)
			if out, err := w.svixClient.Message.CreateWithOptions(
				w.svixAppId, messageIn, options); err != nil {
				w.errChan <- err
				err = fmt.Errorf(
					"svix.Svix.Message.CreateWithOptions: %s: dumpIn: %s", err, dumpIn)
				sharedlogging.GetLogger(ctx).Errorf(err.Error())
			} else {
				w.messageChan <- out
				dumpOut := spew.Sdump(out)
				sharedlogging.GetLogger(ctx).Debug("new webhook sent: ",
					"dumpIn: ", dumpIn, "dumpOut: ", dumpOut)
			}
		} else {
			sharedlogging.GetLogger(ctx).Debug("svix.Svix.EventType.Get: ", "ev: ", ev, "err: ", err)
		}

		if err = w.reader.CommitMessages(ctx, m); err != nil {
			w.errChan <- err
			sharedlogging.GetLogger(ctx).Errorf("kafka.Reader.CommitMessages: %s", err)
			continue
		}
	}
}

func (w *Worker) Stop(ctx context.Context) {
	ch := make(chan struct{})
	w.stopChan <- ch
	sharedlogging.GetLogger(ctx).Debug("worker stopped")
}

func (w *Worker) closeChannels() {
	close(w.stopChan)
	close(w.eventChan)
	close(w.messageChan)
	close(w.errChan)
}
