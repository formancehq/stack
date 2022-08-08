package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks-cloud/internal/storage"
	svix "github.com/svix/svix-webhooks/go"
)

type Event struct {
	Date    time.Time      `json:"date"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

type Worker struct {
	reader     Reader
	store      storage.Store
	svixClient *svix.Svix
	svixAppId  string
}

func NewWorker(reader Reader, store storage.Store, svixClient *svix.Svix, svixAppId string) *Worker {
	return &Worker{
		reader:     reader,
		store:      store,
		svixClient: svixClient,
		svixAppId:  svixAppId,
	}
}

func (e *Worker) Run(ctx context.Context) (fetchedMsgs, sentWebhooks int, err error) {
	sharedlogging.GetLogger(ctx).Infof("starting to read messages")

	for {
		m, err := e.reader.FetchMessage(ctx)
		if err != nil {
			sharedlogging.GetLogger(ctx).Errorf("unable to fetch messages: %s", err)
			if <-ctx.Done(); true {
				sharedlogging.GetLogger(ctx).Infof("context deadline exceeded")
				return fetchedMsgs, sentWebhooks, nil
			}
			continue
		}

		ctx := sharedlogging.ContextWithLogger(ctx,
			sharedlogging.GetLogger(ctx).WithFields(map[string]any{
				"offset": m.Offset,
			}))
		sharedlogging.GetLogger(ctx).WithFields(map[string]any{
			"time":      m.Time,
			"partition": m.Partition,
			"data":      m.Value,
			"headers":   m.Headers,
		}).Debug("New message")

		ev := Event{}
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			sharedlogging.GetLogger(ctx).Errorf("unable to unmarshal message: %s", err)
			continue
		}

		spew.Dump("EVENT", ev)

		if err = e.reader.CommitMessages(ctx, m); err != nil {
			sharedlogging.GetLogger(ctx).Errorf("unable to commit event: %s", err)
			if <-ctx.Done(); true {
				sharedlogging.GetLogger(ctx).Infof("context deadline exceeded")
				return fetchedMsgs, sentWebhooks, nil
			}
			continue
		}

		sharedlogging.GetLogger(ctx).Infof(
			"new message read: %s %s", ev.Date.Format(time.RFC3339), ev.Type)
		fetchedMsgs++

		if _, err := e.svixClient.EventType.Get(ev.Type); err == nil {
			id := uuid.New().String()
			if _, err := e.svixClient.Message.CreateWithOptions(e.svixAppId, &svix.MessageIn{
				EventType: ev.Type,
				EventId:   *svix.NullableString(id),
				Payload:   ev.Payload,
			}, &svix.PostOptions{IdempotencyKey: &id}); err != nil {
				err = fmt.Errorf("unable to send message to %s: %s", e.svixAppId, err)
				sharedlogging.GetLogger(ctx).Errorf(err.Error())
				return fetchedMsgs, sentWebhooks, err
			}
			sharedlogging.GetLogger(ctx).Infof(
				"new webhook sent: %s %s", ev.Date.Format(time.RFC3339), ev.Type)
			sentWebhooks++
		}
	}
}
