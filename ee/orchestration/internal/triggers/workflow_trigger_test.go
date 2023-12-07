package triggers

import (
	"testing"
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestWorkflow(t *testing.T) {
	t.Parallel()

	testSuite := &testsuite.WorkflowTestSuite{}

	req := ProcessEventRequest{
		Event: publish.EventMessage{
			Date: time.Now().Round(time.Second).UTC(),
		},
	}

	trigger := Trigger{
		TriggerData: TriggerData{
			Event:      "NEW_TRANSACTION",
			WorkflowID: "xxx",
		},
		ID: uuid.NewString(),
	}

	env := testSuite.NewTestWorkflowEnvironment()
	env.
		OnActivity(ListTriggersActivity, mock.Anything, req).
		Once().
		Return([]Trigger{trigger}, nil)
	env.
		OnActivity(ProcessEventActivity, mock.Anything, trigger, req).
		Once().
		Return(nil)

	env.ExecuteWorkflow(RunTrigger, req)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
