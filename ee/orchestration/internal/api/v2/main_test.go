package v2

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/testing/docker"
	"github.com/formancehq/go-libs/testing/utils"

	"github.com/formancehq/go-libs/bun/bundebug"

	"go.temporal.io/sdk/worker"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/google/uuid"
	"go.temporal.io/sdk/testsuite"

	"github.com/formancehq/go-libs/bun/bunconnect"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/orchestration/internal/triggers"

	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/testing/platform/pgtesting"
	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func test(t *testing.T, fn func(router *chi.Mux, backend api.Backend, db *bun.DB)) {
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

	taskQueue := uuid.NewString()
	worker := temporalworker.New(logging.Testing(), devServer.Client(), taskQueue,
		[]temporalworker.DefinitionSet{
			workflow.NewWorkflows(false).DefinitionSet(),
			temporalworker.NewDefinitionSet().Append(temporalworker.Definition{
				Name: "NoOp",
				Func: (&stages.NoOp{}).GetWorkflow(),
			}),
		},
		[]temporalworker.DefinitionSet{
			workflow.NewActivities(publish.NoOpPublisher, db).DefinitionSet(),
		},
		worker.Options{},
	)
	require.NoError(t, worker.Start())
	t.Cleanup(worker.Stop)

	require.NoError(t, storage.Migrate(context.Background(), db))
	workflowManager := workflow.NewManager(db, devServer.Client(), taskQueue, false)
	expressionEvaluator := triggers.NewExpressionEvaluator(http.DefaultClient)
	triggersManager := triggers.NewManager(db, expressionEvaluator)
	backend := api.NewDefaultBackend(triggersManager, workflowManager)
	router := newRouter(backend, auth.NewNoAuth(), testing.Verbose())
	fn(router, backend, db)
}

var (
	devServer *testsuite.DevServer
	srv       *pgtesting.PostgresServer
)

func TestMain(m *testing.M) {
	utils.WithTestMain(func(t *utils.TestingTForMain) int {
		srv = pgtesting.CreatePostgresServer(t, docker.NewPool(t, logging.Testing()))

		var err error
		devServer, err = testsuite.StartDevServer(logging.TestingContext(), testsuite.DevServerOptions{})
		if err != nil {
			log.Fatal(err)
		}

		t.Cleanup(func() {
			require.NoError(t, devServer.Stop())
		})

		return m.Run()
	})
}
