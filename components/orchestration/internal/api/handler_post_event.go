package api

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func postEventToWorkflowInstance(backend Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event := workflow.Event{}
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			api.BadRequest(w, "VALIDATION", err)
			return
		}

		if err := backend.PostEvent(r.Context(), instanceID(r), event); err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
