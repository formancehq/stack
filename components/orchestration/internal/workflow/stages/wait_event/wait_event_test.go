package wait_event

import (
	"testing"
	"time"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal/stagestesting"
	"go.temporal.io/sdk/testsuite"
)

func TestWaitEventSchemaValidation(t *testing.T) {
	stagestesting.TestSchemas(t, "wait_event", []stagestesting.SchemaTestCase{
		{
			Data: map[string]any{
				"wait_event": map[string]any{},
			},
			ExpectedResolved:        WaitEvent{},
			ExpectedValidationError: true,
		},
		{
			Name: "valid case",
			Data: map[string]any{
				"event": "test",
			},
			ExpectedResolved: WaitEvent{
				Event: "test",
			},
			ExpectedValidationError: false,
		},
	}...)
}

func TestWaitEvent(t *testing.T) {
	stagestesting.RunWorkflows(t, []stagestesting.WorkflowTestCase[WaitEvent]{
		{
			Stage: WaitEvent{
				Event: "test",
			},
			DelayedCallbacks: []stagestesting.DelayedCallback{{
				Fn: func(environment *testsuite.TestWorkflowEnvironment) func() {
					return func() {
						environment.SignalWorkflow(workflow.EventSignalName, workflow.Event{
							Name: "test",
						})
					}
				},
				Duration: 100 * time.Millisecond,
			}},
			Name: "nominal",
		},
	}...)
}
