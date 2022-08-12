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
}

func NewWorker(reader Reader, store storage.Store, svixClient *svixgo.Svix, svixAppId string) *Worker {
	return &Worker{
		reader:     reader,
		store:      store,
		svixClient: svixClient,
		svixAppId:  svixAppId,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			sharedlogging.GetLogger(ctx).Debug("worker: context done")
			return nil
		default:
		}

		m, err := w.reader.FetchMessage(ctx)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				err := fmt.Errorf("kafka.Reader.FetchMessage: %w", err)
				sharedlogging.GetLogger(ctx).Errorf(err.Error())
				return err
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
			sharedlogging.GetLogger(ctx).Errorf("json.Unmarshal: %s", err)
			continue
		}

		if out, err := w.svixClient.EventType.Get(ev.Type); err == nil {
			sharedlogging.GetLogger(ctx).Debug("svix.Svix.EventType.Get: ",
				"ev.Type: ", ev.Type, "dumpOut: ", spew.Sdump(out))

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
				err = fmt.Errorf("svix.Svix.Message.CreateWithOptions: %s: dumpIn: %s",
					err, dumpIn)
				sharedlogging.GetLogger(ctx).Errorf(err.Error())
				continue
			} else {
				sharedlogging.GetLogger(ctx).Debug("new webhook sent: ",
					"dumpIn: ", dumpIn, "dumpOut: ", spew.Sdump(out))
			}
		} else {
			sharedlogging.GetLogger(ctx).Debug("svix.Svix.EventType.Get: ",
				"ev: ", ev, "err: ", err)
		}

		if err := w.reader.CommitMessages(ctx, m); err != nil {
			err := fmt.Errorf("kafka.Reader.CommitMessages: %w", err)
			sharedlogging.GetLogger(ctx).Errorf(err.Error())
			return err
		}
	}
}
