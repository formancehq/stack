package triggers

import (
	"testing"
	"time"

	worker "go.temporal.io/sdk/worker"

	"github.com/formancehq/go-libs/bun/bundebug"
	"github.com/uptrace/bun"

	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"go.temporal.io/sdk/client"

	"github.com/formancehq/go-libs/publish"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWorkflow(t *testing.T) {
	t.Parallel()

	hooks := make([]bun.QueryHook, 0)
	if testing.Verbose() {
		hooks = append(hooks, bundebug.NewQueryHook())
	}

	database := srv.NewDatabase(t)
	db, err := bunconnect.OpenSQLDB(logging.TestingContext(), bunconnect.ConnectionOptions{
		DatabaseSourceName: database.ConnString(),
	}, hooks...)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})
	require.NoError(t, storage.Migrate(logging.TestingContext(), db))

	taskQueue := uuid.NewString()
	workflowManager := workflow.NewManager(db, devServer.Client(), taskQueue, false)

	worker := temporalworker.New(logging.Testing(), devServer.Client(), taskQueue,
		[]temporalworker.DefinitionSet{
			NewWorkflow(taskQueue, false).DefinitionSet(),
			workflow.NewWorkflows(false).DefinitionSet(),
			temporalworker.NewDefinitionSet().Append(temporalworker.Definition{
				Name: "NoOp",
				Func: (&stages.NoOp{}).GetWorkflow(),
			}),
		},
		[]temporalworker.DefinitionSet{
			workflow.NewActivities(publish.NoOpPublisher, db).DefinitionSet(),
			NewActivities(db, workflowManager, NewDefaultExpressionEvaluator(), publish.NoOpPublisher).DefinitionSet(),
		},
		worker.Options{},
	)
	require.NoError(t, worker.Start())
	t.Cleanup(worker.Stop)

	req := ProcessEventRequest{
		Event: publish.EventMessage{
			Type: "NEW_TRANSACTION",
			Date: time.Now().Round(time.Second).UTC(),
		},
	}

	workflow := workflow.New(workflow.Config{
		Stages: []workflow.RawStage{{
			"noop": map[string]any{},
		}},
	})
	_, err = db.
		NewInsert().
		Model(&workflow).
		Exec(logging.TestingContext())
	require.NoError(t, err)

	trigger := Trigger{
		TriggerData: TriggerData{
			Event:      "NEW_TRANSACTION",
			WorkflowID: workflow.ID,
		},
		ID: uuid.NewString(),
	}
	_, err = db.NewInsert().Model(&trigger).Exec(logging.TestingContext())
	require.NoError(t, err)

	ret, err := devServer.Client().
		ExecuteWorkflow(logging.TestingContext(), client.StartWorkflowOptions{
			TaskQueue: taskQueue,
		}, RunTrigger, req)
	require.NoError(t, err)
	require.NoError(t, ret.Get(logging.TestingContext(), nil))
}
