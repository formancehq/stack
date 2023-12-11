package v1

import (
	"net/http"

	"github.com/formancehq/orchestration/internal/workflow"

	api "github.com/formancehq/orchestration/internal/api"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
)

func listInstances(backend api.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		runs, err := backend.ListInstances(r.Context(), workflow.ListInstancesQuery{
			Options: workflow.ListInstancesOptions{
				WorkflowID: r.URL.Query().Get("workflowID"),
				Running:    sharedapi.QueryParamBool(r, "running"),
			},
		})
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}
		sharedapi.Ok(w, runs.Data)
	}
}
