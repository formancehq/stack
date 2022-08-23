package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/pkg/svix"
	svixgo "github.com/svix/svix-webhooks/go"
)

var (
	ErrDateZero     = errors.New("date cannot be zero")
	ErrTypeEmpty    = errors.New("type cannot be empty")
	ErrPayloadEmpty = errors.New("payload cannot be empty")
)

type Event struct {
	Date    time.Time      `json:"date"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (e Event) Validate() error {
	if e.Date.IsZero() {
		return ErrDateZero
	}

	if e.Type == "" {
		return ErrTypeEmpty
	}

	if len(e.Payload) == 0 {
		return ErrPayloadEmpty
	}

	return nil
}

func processEventMessage(ctx context.Context, msgValue []byte, svixApp svix.App) error {
	ev := Event{}
	if err := json.Unmarshal(msgValue, &ev); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	if err := ev.Validate(); err != nil {
		return fmt.Errorf("kafka.Event.Validate: %w", err)
	}

	id := uuid.NewString()
	messageIn := &svixgo.MessageIn{
		EventType: ev.Type,
		EventId:   *svixgo.NullableString(id),
		Payload:   ev.Payload,
	}
	options := &svixgo.PostOptions{IdempotencyKey: &id}
	dumpIn := spew.Sdump(
		"svix appId: ", svixApp.AppId,
		"svix.MessageIn: ", messageIn,
		"svix.PostOptions: ", options)

	if out, err := svixApp.Client.Message.CreateWithOptions(
		svixApp.AppId, messageIn, options); err != nil {
		return fmt.Errorf("svix.Svix.Message.CreateWithOptions: %w: dumpIn: %s",
			err, dumpIn)
	} else {
		sharedlogging.GetLogger(ctx).Debug("new webhook sent: ",
			"dumpIn: ", dumpIn, "dumpOut: ", spew.Sdump(out))
	}

	return nil
}
