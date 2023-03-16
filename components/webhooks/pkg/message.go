package webhooks

import (
	"time"

	"github.com/formancehq/stack/libs/events"
)

type Response struct {
	Date    time.Time `json:"date"`
	App     string    `json:"app"`
	Payload any       `json:"payload"`
}

func BuildResponseFromMessage(msg *events.Event) *Response {
	return &Response{
		Date:    msg.CreatedAt.AsTime(),
		App:     msg.App,
		Payload: msg.Event,
	}
}
