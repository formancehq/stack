package api

import (
	"context"
	"fmt"
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

func TestGetInstance(t *testing.T) {
	test(t, func(router *chi.Mux, m Backend, db *bun.DB) {
		w, err := m.Create(context.TODO(), workflow.Config{
			Stages: []workflow.RawStage{},
		})
		require.NoError(t, err)

		instance, err := m.RunWorkflow(context.TODO(), w.ID, map[string]string{})
		require.NoError(t, err)

		now := time.Now().Round(time.Nanosecond)
		for i := 0; i < 10; i++ {
			timestamp := now.Add(time.Second)
			_, err := db.NewInsert().Model(&workflow.Stage{
				Number:       i,
				InstanceID:   instance.ID,
				StartedAt:    now,
				TerminatedAt: &timestamp,
			}).Exec(context.TODO())
			require.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/instances/%s", instance.ID), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		var retrievedInstance workflow.Instance
		apitesting.ReadResponse(t, rec, &retrievedInstance)
		require.Len(t, retrievedInstance.Statuses, 10)
	})
}
