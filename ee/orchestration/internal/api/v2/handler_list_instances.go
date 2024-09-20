package v2

import (
	"net/http"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/orchestration/internal/workflow"

	api "github.com/formancehq/orchestration/internal/api"

	sharedapi "github.com/formancehq/go-libs/api"
)

func listInstances(backend api.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		q, err := bunpaginate.Extract[workflow.ListInstancesQuery](r, func() (*workflow.ListInstancesQuery, error) {
			pageSize, err := bunpaginate.GetPageSize(r)
			if err != nil {
				return nil, err
			}
			return &workflow.ListInstancesQuery{
				PageSize: pageSize,
				Options: workflow.ListInstancesOptions{
					WorkflowID: r.URL.Query().Get("workflowID"),
					Running:    sharedapi.QueryParamBool(r, "running"),
				},
			}, nil
		})
		if err != nil {
			sharedapi.BadRequest(w, "VALIDATION", err)
			return
		}

		runs, err := backend.ListInstances(r.Context(), *q)
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}
		sharedapi.RenderCursor(w, *runs)
	}
}
