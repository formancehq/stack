package v3

import (
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
)

func workflowsInstancesList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v3_workflowsInstancesList")
		defer span.End()

		connectorID, err := models.ConnectorIDFromString(connectorID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		scheduleID := scheduleID(r)

		query, err := bunpaginate.Extract[storage.ListInstancesQuery](r, func() (*storage.ListInstancesQuery, error) {
			pageSize, err := bunpaginate.GetPageSize(r)
			if err != nil {
				return nil, err
			}

			options := pointer.For(bunpaginate.NewPaginatedQueryOptions(storage.InstanceQuery{}).WithPageSize(pageSize))
			options = pointer.For(options.WithQueryBuilder(
				query.And(
					query.Match("connector_id", connectorID),
					query.Match("schedule_id", scheduleID),
				),
			))

			return pointer.For(storage.NewListInstancesQuery(*options)), nil
		})
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := backend.WorkflowsInstancesList(ctx, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.RenderCursor(w, *cursor)
	}
}
