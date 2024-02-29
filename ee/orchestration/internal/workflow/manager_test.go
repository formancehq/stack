package workflow

import (
	"testing"
	"time"

	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	database := pgtesting.NewPostgresDatabase(t)
	db, err := bunconnect.OpenSQLDB(logging.TestingContext(), bunconnect.ConnectionOptions{
		DatabaseSourceName: database.ConnString(),
		Debug:              testing.Verbose(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})
	require.NoError(t, storage.Migrate(logging.TestingContext(), db))

	taskQueue := uuid.NewString()
	worker := temporalworker.New(logging.Testing(), devServer.Client(), taskQueue,
		[]any{NewWorkflows(), (&stages.NoOp{}).GetWorkflow()},
		[]any{NewActivities(publish.NoOpPublisher, db)},
	)
	require.NoError(t, worker.Start())
	t.Cleanup(worker.Stop)

	manager := NewManager(db, devServer.Client(), taskQueue)

	config := Config{
		Stages: []RawStage{
			{
				"noop": map[string]any{},
			},
		},
	}
	w, err := manager.Create(logging.TestingContext(), config)
	require.NoError(t, err)

	i, err := manager.RunWorkflow(logging.TestingContext(), w.ID, map[string]string{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		updatedInstance, err := manager.GetInstance(logging.TestingContext(), i.ID)
		require.NoError(t, err)
		return len(updatedInstance.Statuses) == 1
	}, 2*time.Second, 100*time.Millisecond)
}

/**
register activity InsertNewInstance
register activity InsertNewStage
register activity SendWorkflowTerminationEvent
register activity UpdateInstance
register activity UpdateStage
*/
