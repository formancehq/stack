package events

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

const (
	SucceededWorkflow = "SUCCEEDED_WORKFLOW"
	FailedWorkflow    = "FAILED_WORKFLOW"
	SucceededTrigger  = "SUCCEEDED_TRIGGER"
	FailedTrigger     = "FAILED_TRIGGER"
)

type SucceededWorkflowPayload struct {
	ID         string `json:"id"`
	InstanceID string `json:"instanceID"`
}

type FailedWorkflowPayload struct {
	ID         string `json:"id"`
	InstanceID string `json:"instanceID"`
	Error      string `json:"error"`
}

type SucceededTriggerPayload struct {
	ID      string `json:"id"`
	EventID string `json:"eventID"`
}

type FailedTriggerPayload struct {
	ID      string `json:"id"`
	EventID string `json:"eventID"`
	Error   string `json:"error"`
}

func NewMessage(ctx context.Context, mtype string, payload any) *message.Message {
	return publish.NewMessage(ctx, publish.EventMessage{
		Date:    time.Now(),
		App:     "orchestration",
		Version: "v2",
		Type:    mtype,
		Payload: payload,
	})
}
