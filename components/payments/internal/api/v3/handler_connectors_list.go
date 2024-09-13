package v3

import (
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/payments/internal/storage"
)

func connectorsList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v3_connectorsList")
		defer span.End()

		query, err := bunpaginate.Extract[storage.ListConnectorsQuery](r, func() (*storage.ListConnectorsQuery, error) {
			options, err := getPagination(r, storage.ConnectorQuery{})
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListConnectorsQuery(*options)), nil
		})
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		connectors, err := backend.ConnectorsList(ctx, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.RenderCursor(w, *connectors)
	}
}
