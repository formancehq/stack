package v1

import (
	"errors"
	"net/http"

	api2 "github.com/formancehq/orchestration/internal/api"
	"github.com/go-playground/validator/v10"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api"
)

var (
	ErrEmptyID = errors.New("ID is empty")
)

func deleteWorkflow(backend api2.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := workflowID(r)

		err := validator.New().Var(id, "required,uuid")
		if err != nil {
			api.BadRequest(w, "VALIDATION", err)
			return
		}

		err = backend.DeleteWorkflow(r.Context(), workflowID(r))

		if errors.Is(err, workflow.ErrWorkflowNotFound) {
			api.NotFound(w, err)
			return
		}

		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
