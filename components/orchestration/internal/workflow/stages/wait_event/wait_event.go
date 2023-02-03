package wait_event

import (
	"github.com/formancehq/orchestration/internal/workflow/stages"
)

type WaitEvent struct {
	Event string `json:"event" validate:"required"`
}

func (w WaitEvent) GetWorkflow() any {
	return RunWaitEvent
}

func init() {
	stages.Register("wait_event", WaitEvent{})
}
