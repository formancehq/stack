package v2

import (
	"net/http"
	"strconv"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/go-chi/chi/v5"

	"github.com/formancehq/orchestration/internal/workflow"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pkg/errors"
)

func readStageHistory(backend api.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stageNumberAsString := chi.URLParam(r, "number")
		stage, err := strconv.ParseInt(stageNumberAsString, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		workflows, err := backend.ReadStageHistory(r.Context(), instanceID(r), int(stage))
		if err != nil {
			switch {
			case errors.Is(err, workflow.ErrInstanceNotFound):
				sharedapi.NotFound(w)
			default:
				sharedapi.InternalServerError(w, r, err)
			}
			return
		}

		sharedapi.Ok(w, workflows)
	}
}
