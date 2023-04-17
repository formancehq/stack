package api

import (
	"net/http"
	"strconv"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func readStageHistory(m *workflow.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stageNumberAsString := chi.URLParam(r, "number")
		stage, err := strconv.ParseInt(stageNumberAsString, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		workflows, err := m.ReadStageHistory(r.Context(), instanceID(r), int(stage))
		if err != nil {
			switch {
			case errors.Is(err, workflow.ErrInstanceNotFound):
				api.NotFound(w)
			default:
				api.InternalServerError(w, r, err)
			}
			return
		}

		api.Ok(w, workflows)
	}
}
