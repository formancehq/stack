package workflow

import (
	"context"
	"os"
	"testing"

	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/testsuite"
)

type mockTemporalClient struct {
	client.Client
	env       *testsuite.TestWorkflowEnvironment
	workflows *Workflows
}

func (c *mockTemporalClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	c.env.ExecuteWorkflow(c.workflows.Run, args...)
	return nil, nil
}

func TestConfig(t *testing.T) {
	t.Parallel()

	database := pgtesting.NewPostgresDatabase(t)
	db := storage.LoadDB(database.ConnString(), testing.Verbose(), os.Stdout)
	t.Cleanup(func() {
		_ = db.Close()
	})

	require.NoError(t, storage.Migrate(context.Background(), db))
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	workflows := NewWorkflows(db)
	env.RegisterWorkflow(stages.RunNoOp)
	env.RegisterWorkflow(workflows.Run)
	mockClient := &mockTemporalClient{env: env, workflows: workflows}
	manager := NewManager(db, mockClient, "default")

	config := Config{
		Stages: []RawStage{
			{
				"noop": map[string]any{},
			},
		},
	}
	w, err := manager.Create(context.TODO(), config)
	require.NoError(t, err)

	i, err := manager.RunWorkflow(context.TODO(), w.ID, map[string]string{})
	require.NoError(t, err)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	updatedInstance, err := manager.GetInstance(context.TODO(), i.ID)
	require.NoError(t, err)
	require.Len(t, updatedInstance.Statuses, 1)
}
