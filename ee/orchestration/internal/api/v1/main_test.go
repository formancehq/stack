package v1

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bundebug"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/temporal"
	"github.com/formancehq/stack/libs/go-libs/testing/docker"
	"github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
	"github.com/formancehq/stack/libs/go-libs/testing/utils"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/worker"
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
	worker := temporal.New(logging.Testing(), devServer.Client(), taskQueue,
		[]temporal.DefinitionSet{
			workflow.NewWorkflows(false).DefinitionSet(),
			temporal.NewDefinitionSet().Append(temporal.Definition{
				Name: "NoOp",
				Func: (&stages.NoOp{}).GetWorkflow(),
			}),
		},
		[]temporal.DefinitionSet{
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
	srv       *pgtesting.PostgresServer
	devServer *testsuite.DevServer
)

func TestMain(m *testing.M) {
	flag.Parse()

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
