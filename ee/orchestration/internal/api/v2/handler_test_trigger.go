package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/orchestration/internal/api"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/go-chi/chi/v5"
)

func testTrigger(backend api.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := make(map[string]any)
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		o, err := backend.TestTrigger(r.Context(), chi.URLParam(r, "triggerID"), data)
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		sharedapi.Ok(w, o)
	}
}
