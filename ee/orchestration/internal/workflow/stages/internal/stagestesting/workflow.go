package stagestesting

import (
	"testing"
	"time"

	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

type MockedActivity struct {
	Activity any
	Args     []any
	Returns  []any
}

type DelayedCallback struct {
	Fn       func(environment *testsuite.TestWorkflowEnvironment) func()
	Duration time.Duration
}

type WorkflowTestCase[T stages.Stage] struct {
	Stage            T
	MockedActivities []MockedActivity
	DelayedCallbacks []DelayedCallback
	Name             string
}

func RunWorkflowTest[T stages.Stage](t *testing.T, testCase WorkflowTestCase[T]) {
	t.Run(testCase.Name, func(t *testing.T) {
		t.Parallel()

		testSuite := &testsuite.WorkflowTestSuite{}

		env := testSuite.NewTestWorkflowEnvironment()
		for _, ma := range testCase.MockedActivities {
			env.OnActivity(ma.Activity, ma.Args...).Return(ma.Returns...)
		}
		for _, callback := range testCase.DelayedCallbacks {
			env.RegisterDelayedCallback(callback.Fn(env), callback.Duration)
		}

		var stage T
		env.ExecuteWorkflow(stage.GetWorkflow(), testCase.Stage)
		require.True(t, env.IsWorkflowCompleted())
		require.NoError(t, env.GetWorkflowError())
	})
}

func RunWorkflows[T stages.Stage](t *testing.T, testCases ...WorkflowTestCase[T]) {
	t.Parallel()

	for _, testCase := range testCases {
		RunWorkflowTest(t, testCase)
	}
}
