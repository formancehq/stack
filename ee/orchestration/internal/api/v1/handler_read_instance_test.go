package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sharedapi "github.com/formancehq/stack/libs/go-libs/testing/api"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestGetInstance(t *testing.T) {
	test(t, func(router *chi.Mux, m api.Backend, db *bun.DB) {
		w, err := m.Create(logging.TestingContext(), workflow.Config{
			Stages: []workflow.RawStage{},
		})
		require.NoError(t, err)

		instance := workflow.NewInstance(uuid.NewString(), w.ID)
		_, err = db.NewInsert().
			Model(&instance).
			Exec(logging.TestingContext())
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
		sharedapi.ReadResponse(t, rec, &retrievedInstance)
		require.Len(t, retrievedInstance.Statuses, 10)
	})
}
