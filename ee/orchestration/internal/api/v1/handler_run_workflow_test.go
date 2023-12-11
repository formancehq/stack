package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/go-chi/chi/v5"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api/apitesting"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestRunWorkflow(t *testing.T) {
	test(t, func(router *chi.Mux, m api.Backend, db *bun.DB) {
		w, err := m.Create(context.TODO(), workflow.Config{
			Stages: []workflow.RawStage{},
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/workflows/%s/instances", w.ID), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	})
}

func TestRunWorkflowWaitEvent(t *testing.T) {
	test(t, func(router *chi.Mux, m api.Backend, db *bun.DB) {
		w, err := m.Create(context.TODO(), workflow.Config{
			Stages: []workflow.RawStage{
				map[string]map[string]any{
					"noop": {},
				},
			},
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/workflows/%s/instances?wait=true", w.ID), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)
		instance := &workflow.Instance{}
		apitesting.ReadResponse(t, rec, instance)
	})
}
