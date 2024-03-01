package v2

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/operatorservice/v1"

	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.temporal.io/sdk/testsuite"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/orchestration/internal/triggers"

	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func test(t *testing.T, fn func(router *chi.Mux, backend api.Backend, db *bun.DB)) {
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

	taskQueue := uuid.NewString()
	worker := temporalworker.New(logging.Testing(), devServer.Client(), taskQueue,
		[]any{workflow.NewWorkflows(), (&stages.NoOp{}).GetWorkflow()},
		[]any{workflow.NewActivities(publish.NoOpPublisher, db)},
	)
	require.NoError(t, worker.Start())
	t.Cleanup(worker.Stop)

	require.NoError(t, storage.Migrate(context.Background(), db))
	workflowManager := workflow.NewManager(db, devServer.Client(), taskQueue)
	expressionEvaluator := triggers.NewExpressionEvaluator(http.DefaultClient)
	triggersManager := triggers.NewManager(db, expressionEvaluator)
	backend := api.NewDefaultBackend(triggersManager, workflowManager)
	router := newRouter(backend, auth.NewNoAuth())
	fn(router, backend, db)
}

var (
	devServer *testsuite.DevServer
)

func TestMain(m *testing.M) {
	flag.Parse()

	if err := pgtesting.CreatePostgresServer(); err != nil {
		log.Fatal(err)
	}

	var err error
	devServer, err = testsuite.StartDevServer(logging.TestingContext(), testsuite.DevServerOptions{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = devServer.Client().OperatorService().AddSearchAttributes(logging.TestingContext(), &operatorservice.AddSearchAttributesRequest{
		SearchAttributes: map[string]enums.IndexedValueType{
			workflow.SearchAttributeWorkflowID: enums.INDEXED_VALUE_TYPE_TEXT,
			triggers.SearchAttributeTriggerID:  enums.INDEXED_VALUE_TYPE_TEXT,
		},
		Namespace: "default",
	})
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	if err := devServer.Stop(); err != nil {
		log.Println("unable to stop temporal server", err)
	}
	if err := pgtesting.DestroyPostgresServer(); err != nil {
		log.Println("unable to stop postgres server", err)
	}
	os.Exit(code)
}
