package v1

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/go-chi/chi/v5"

	"github.com/formancehq/orchestration/internal/triggers"

	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/client"
)

type mockedRun struct {
	client.WorkflowRun
	runID string
}

func (m mockedRun) GetRunID() string {
	return m.runID
}

func (m mockedRun) Get(ctx context.Context, valuePtr interface{}) error {
	return nil
}

type mockedClient struct {
	client.Client
	db        *bun.DB
	t         *testing.T
	workflows map[string]mockedRun
}

func (c *mockedClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, w interface{}, args ...interface{}) (client.WorkflowRun, error) {
	input := args[0].(workflow.Input)
	for ind := range input.Workflow.Config.Stages {
		timestamp := time.Now()
		_, err := c.db.NewInsert().Model(&workflow.Stage{
			Number:       ind,
			InstanceID:   options.ID,
			StartedAt:    time.Now(),
			TerminatedAt: &timestamp,
		}).Exec(context.Background())
		require.NoError(c.t, err)
	}
	r := mockedRun{
		runID: options.ID,
	}
	c.workflows[options.ID] = r
	return r, nil
}

func (c *mockedClient) GetWorkflow(ctx context.Context, workflowID string, runID string) client.WorkflowRun {
	return c.workflows[workflowID]
}

func newMockedClient(t *testing.T, db *bun.DB) *mockedClient {
	return &mockedClient{
		db:        db,
		t:         t,
		workflows: map[string]mockedRun{},
	}
}

func test(t *testing.T, fn func(router *chi.Mux, backend api.Backend, db *bun.DB)) {
	t.Parallel()

	database := pgtesting.NewPostgresDatabase(t)
	db := storage.LoadDB(database.ConnString(), testing.Verbose(), os.Stdout)
	t.Cleanup(func() {
		_ = db.Close()
	})

	require.NoError(t, storage.Migrate(context.Background(), db))
	workflowManager := workflow.NewManager(db, newMockedClient(t, db), "default")
	expressionEvaluator := triggers.NewExpressionEvaluator(http.DefaultClient)
	triggersManager := triggers.NewManager(db, expressionEvaluator)
	backend := api.NewDefaultBackend(triggersManager, workflowManager)
	router := newRouter(backend, auth.NewNoAuth())
	fn(router, backend, db)
}

func TestMain(m *testing.M) {
	flag.Parse()

	if err := pgtesting.CreatePostgresServer(); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	if err := pgtesting.DestroyPostgresServer(); err != nil {
		log.Println("unable to stop postgres server", err)
	}
	os.Exit(code)
}
