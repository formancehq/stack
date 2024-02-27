package workflow

import (
	"testing"

	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestActivities(t *testing.T) {
	publisher := publish.InMemory()
	activities := NewActivities(publisher)

	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(activities.SendWorkflowTerminationEvent)
	_, err := env.ExecuteActivity(SendWorkflowTerminationEventActivity, NewInstance("xxx"))
	require.NoError(t, err)
	require.NotEmpty(t, publisher.AllMessages())
}
