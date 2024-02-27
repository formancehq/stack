package workflow

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/orchestration/pkg/events"
)

type Activities struct {
	publisher message.Publisher
}

func (a Activities) SendWorkflowTerminationEvent(ctx context.Context, instance Instance) error {
	if instance.Error == "" {
		return a.publisher.Publish(events.SucceededWorkflow,
			events.NewMessage(ctx, events.SucceededWorkflow, events.SucceededWorkflowPayload{
				ID:         instance.WorkflowID,
				InstanceID: instance.ID,
			}))
	} else {
		return a.publisher.Publish(events.FailedWorkflow,
			events.NewMessage(ctx, events.FailedWorkflow, events.FailedWorkflowPayload{
				ID:         instance.WorkflowID,
				InstanceID: instance.ID,
				Error:      instance.Error,
			}))
	}
}

var SendWorkflowTerminationEventActivity = (&Activities{}).SendWorkflowTerminationEvent

func NewActivities(publisher message.Publisher) Activities {
	return Activities{
		publisher: publisher,
	}
}
