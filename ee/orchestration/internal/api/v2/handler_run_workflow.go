package v2

import (
	"encoding/json"
	"net/http"
	"strings"

	api2 "github.com/formancehq/orchestration/internal/api"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func runWorkflow(backend api2.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := make(map[string]string)
		if r.ContentLength > 0 {
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				api.BadRequest(w, "VALIDATION", err)
				return
			}
		}
		instance, err := backend.RunWorkflow(r.Context(), workflowID(r), input)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		if wait := strings.ToLower(r.URL.Query().Get("wait")); wait == "true" || wait == "1" {
			ret := struct {
				*workflow.Instance
				Error string `json:"error,omitempty"`
			}{
				Instance: instance,
			}
			if err := backend.Wait(r.Context(), instance.ID); err != nil {
				ret.Error = err.Error()
			}
			ret.Instance, err = backend.GetInstance(r.Context(), instance.ID)
			if err != nil {
				panic(err)
			}

			api.Created(w, ret)
			return
		}

		api.Created(w, instance)
	}
}
