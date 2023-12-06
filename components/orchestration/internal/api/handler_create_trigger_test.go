package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestCreateTrigger(t *testing.T) {
	test(t, func(router *chi.Mux, m Backend, db *bun.DB) {

		w, err := m.Create(context.Background(), workflow.Config{})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/triggers",
			bytes.NewBufferString(fmt.Sprintf(`{"workflowID": "%s", "event": "SAVED_PAYMENT"}`, w.ID)))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	})
}
