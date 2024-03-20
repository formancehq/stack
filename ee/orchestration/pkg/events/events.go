package events

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

const (
	TopicOrchestration = "orchestration"

	SucceededWorkflow      = "SUCCEEDED_WORKFLOW"
	SucceededWorkflowStage = "SUCCEEDED_WORKFLOW_STAGE"
	StartedWorkflow        = "STARTED_WORKFLOW"
	StartedWorkflowStage   = "STARTED_WORKFLOW_STAGE"
	FailedWorkflow         = "FAILED_WORKFLOW"
	FailedWorkflowStage    = "FAILED_WORKFLOW_STAGE"
	SucceededTrigger       = "SUCCEEDED_TRIGGER"
	FailedTrigger          = "FAILED_TRIGGER"
)

type SucceededWorkflowPayload struct {
	ID         string `json:"id"`
	InstanceID string `json:"instanceID"`
}

type StartedWorkflowPayload struct {
	ID         string `json:"id"`
	InstanceID string `json:"instanceID"`
}

type StartedWorkflowStagePayload struct {
	ID         string `json:"id"`
	InstanceID string `json:"instanceID"`
	Number     int    `json:"number"`
}

type SucceededWorkflowStagePayload struct {
	ID         string `json:"id"`
	InstanceID string `json:"instanceID"`
	Number     int    `json:"number"`
}

type FailedWorkflowPayload struct {
	ID         string `json:"id"`
	InstanceID string `json:"instanceID"`
	Error      string `json:"error"`
}

type FailedWorkflowStagePayload struct {
	ID         string `json:"id"`
	InstanceID string `json:"instanceID"`
	Number     int    `json:"number"`
	Error      string `json:"error"`
}

type SucceededTriggerPayload struct {
	ID                 string `json:"id"`
	TriggerID          string `json:"triggerID"`
	WorkflowInstanceID string `json:"workflowInstanceID"`
}

type FailedTriggerPayload struct {
	ID        string `json:"id"`
	TriggerID string `json:"triggerID"`
	Error     string `json:"error"`
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
