package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/go-libs/logging"

	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/go-libs/testing/api"

	"github.com/google/uuid"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestListInstances(t *testing.T) {
	ctx := logging.TestingContext()

	test(t, func(router *chi.Mux, m api.Backend, db *bun.DB) {
		// Create a workflow with 10 instances
		w := workflow.New(workflow.Config{})
		_, err := db.NewInsert().Model(&w).Exec(ctx)
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			instance := workflow.NewInstance(uuid.NewString(), w.ID)
			if i > 5 {
				instance.SetTerminated(time.Now())
			}
			_, err := db.NewInsert().Model(&instance).Exec(ctx)
			require.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/instances", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		instances := make([]workflow.Instance, 0)
		sharedapi.ReadResponse(t, rec, &instances)
		require.Len(t, instances, 10)

		// Retrieve only running instances
		req = httptest.NewRequest(http.MethodGet, "/instances?running=true", nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		sharedapi.ReadResponse(t, rec, &instances)
		require.Len(t, instances, 6)

		// Delete the workflow
		req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/workflows/%s/", w.ID), nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)

		// Try to retrieve instances for all workflows
		req = httptest.NewRequest(http.MethodGet, "/instances", nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		instances = make([]workflow.Instance, 0)
		sharedapi.ReadResponse(t, rec, &instances)
		require.Len(t, instances, 0)

		// Try to retrieve instances for the deleted workflow
		req = httptest.NewRequest(http.MethodGet, "/instances?workflowID="+w.ID, nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		instances = make([]workflow.Instance, 0)
		sharedapi.ReadResponse(t, rec, &instances)
		require.Len(t, instances, 0)
	})
}
