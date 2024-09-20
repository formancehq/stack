package workflow

import (
	"testing"
	"time"

	"github.com/formancehq/go-libs/bun/bundebug"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/worker"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/google/uuid"

	"github.com/formancehq/go-libs/bun/bunconnect"

	"github.com/formancehq/orchestration/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
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
	worker := temporalworker.New(logging.Testing(), devServer.Client(), taskQueue,
		[]temporalworker.DefinitionSet{
			NewWorkflows(false).DefinitionSet(),
			temporalworker.NewDefinitionSet().Append(temporalworker.Definition{
				Name: "NoOp",
				Func: (&stages.NoOp{}).GetWorkflow(),
			}),
		},
		[]temporalworker.DefinitionSet{
			NewActivities(publish.NoOpPublisher, db).DefinitionSet(),
		},
		worker.Options{},
	)
	require.NoError(t, worker.Start())
	t.Cleanup(worker.Stop)

	manager := NewManager(db, devServer.Client(), taskQueue, false)

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
