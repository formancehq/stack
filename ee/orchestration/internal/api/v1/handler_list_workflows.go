package v1

import (
	"net/http"

	"github.com/formancehq/go-libs/bun/bunpaginate"

	api2 "github.com/formancehq/orchestration/internal/api"

	"github.com/formancehq/go-libs/api"
)

func listWorkflows(backend api2.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		workflows, err := backend.ListWorkflows(r.Context(), bunpaginate.OffsetPaginatedQuery[any]{})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.Ok(w, workflows.Data)
	}
}
