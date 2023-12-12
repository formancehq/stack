package v2

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	api2 "github.com/formancehq/orchestration/internal/api"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
)

func listWorkflows(backend api2.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		query, err := bunpaginate.Extract[bunpaginate.OffsetPaginatedQuery[any]](r, func() (*bunpaginate.OffsetPaginatedQuery[any], error) {
			pageSize, err := bunpaginate.GetPageSize(r)
			if err != nil {
				return nil, err
			}
			return &bunpaginate.OffsetPaginatedQuery[any]{
				PageSize: pageSize,
			}, nil
		})
		if err != nil {
			sharedapi.BadRequest(w, "VALIDATION", err)
			return
		}

		workflows, err := backend.ListWorkflows(r.Context(), *query)
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		sharedapi.RenderCursor(w, *workflows)
	}
}
