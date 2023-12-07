package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api/apitesting"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/uptrace/bun"
)

func TestDeleteWorkflow(t *testing.T) {
	test(t, func(router *chi.Mux, m Backend, db *bun.DB) {
		// Create a workflow
		req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewBufferString(`{"stages": []}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)

		workflow := workflow.Workflow{}
		apitesting.ReadResponse(t, rec, &workflow)

		require.NotEmpty(t, workflow.ID)

		// 	Delete the workflow
		req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/workflows/%s/", workflow.ID), nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
	})
}
