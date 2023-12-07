package wait_event

import (
	internalWorkflow "github.com/formancehq/orchestration/internal/workflow"
	"go.temporal.io/sdk/workflow"
)

func RunWaitEvent(ctx workflow.Context, waitEvent WaitEvent) error {
	channel := workflow.GetSignalChannel(ctx, internalWorkflow.EventSignalName)
	return workflow.Await(ctx, func() bool {
		var signal internalWorkflow.Event
		ok := channel.ReceiveAsync(&signal)
		if !ok {
			return false
		}
		if signal.Name != waitEvent.Event {
			workflow.GetLogger(ctx).Debug("receive unexpected event", "event", signal.Name)
			return false
		}
		return true
	})
}
