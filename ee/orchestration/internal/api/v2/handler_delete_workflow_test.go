package v2

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/go-libs/testing/api"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/stretchr/testify/require"

	"github.com/uptrace/bun"
)

func TestDeleteWorkflow(t *testing.T) {
	test(t, func(router *chi.Mux, m api.Backend, db *bun.DB) {
		// Create a workflow
		req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewBufferString(`{"stages": []}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)

		workflow := workflow.Workflow{}
		sharedapi.ReadResponse(t, rec, &workflow)

		require.NotEmpty(t, workflow.ID)

		// 	Delete the workflow
		req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/workflows/%s/", workflow.ID), nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
	})
}
