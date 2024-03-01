package workflow

import (
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"

	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestActivities(t *testing.T) {

	database := pgtesting.NewPostgresDatabase(t)
	db, err := bunconnect.OpenSQLDB(logging.TestingContext(), bunconnect.ConnectionOptions{
		DatabaseSourceName: database.ConnString(),
		Debug:              testing.Verbose(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	publisher := publish.InMemory()
	activities := NewActivities(publisher, db)

	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(activities.SendWorkflowTerminationEvent)
	_, err = env.ExecuteActivity(SendWorkflowTerminationEventActivity, NewInstance("vvv", "xxx"))
	require.NoError(t, err)
	require.NotEmpty(t, publisher.AllMessages())
}
