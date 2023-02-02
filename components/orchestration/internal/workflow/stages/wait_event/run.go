package wait_event

import (
	internalWorkflow "github.com/formancehq/orchestration/internal/workflow"
	"go.temporal.io/sdk/workflow"
)

func RunWaitEvent(ctx workflow.Context, waitEvent WaitEvent) error {
	var signal internalWorkflow.Event
	channel := workflow.GetSignalChannel(ctx, internalWorkflow.EventSignalName)
	for {
		_ = channel.Receive(ctx, &signal)
		if signal.Name != waitEvent.Event {
			workflow.GetLogger(ctx).Debug("receive unexpected event", "event", signal.Name)
			continue
		}
		break
	}

	return nil
}
