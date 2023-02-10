package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api/apitesting"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestListInstances(t *testing.T) {
	test(t, func(router *chi.Mux, m *workflow.Manager, db *bun.DB) {

		w := workflow.New(workflow.Config{})
		_, err := db.NewInsert().Model(&w).Exec(context.TODO())
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			instance := workflow.NewInstance(w.ID)
			if i > 5 {
				instance.SetTerminated(time.Now())
			}
			_, err := db.NewInsert().Model(&instance).Exec(context.TODO())
			require.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/instances", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		instances := make([]workflow.Instance, 0)
		apitesting.ReadResponse(t, rec, &instances)
		require.Len(t, instances, 10)

		req = httptest.NewRequest(http.MethodGet, "/instances?running=true", nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		apitesting.ReadResponse(t, rec, &instances)
		require.Len(t, instances, 6)
	})
}
