package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestCreateWorkflow(t *testing.T) {
	test(t, func(router *chi.Mux, m Backend, db *bun.DB) {
		req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewBufferString(`{"stages": []}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	})
}
