package events

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	eventsdef "github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

type Events struct {
	publisher message.Publisher

	stackURL string
}

func New(p message.Publisher, stackURL string) *Events {
	return &Events{
		publisher: p,
		stackURL:  stackURL,
	}
}

func (e *Events) Publish(ctx context.Context, em publish.EventMessage) error {
	return e.publisher.Publish(eventsdef.TopicPayments,
		publish.NewMessage(ctx, em))
}
