package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

// NOTE: in order to maintain previous version compatibility, we need to keep the
// same response structure as the previous version of the API
type connectorsListElement struct {
	Provider    string `json:"provider"`
	ConnectorID string `json:"connectorID"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
}

func connectorsList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_connectorsList")
		defer span.End()

		connectors, err := backend.ConnectorsList(
			ctx,
			storage.NewListConnectorsQuery(
				bunpaginate.NewPaginatedQueryOptions(storage.ConnectorQuery{}).
					// NOTE: previous version of payments did not have pagination, so
					// fetch everything and return it all
					WithPageSize(1000),
			),
		)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]*connectorsListElement, len(connectors.Data))
		for i := range connectors.Data {
			data[i] = &connectorsListElement{
				Provider:    connectors.Data[i].Provider,
				ConnectorID: connectors.Data[i].ID.String(),
				Name:        connectors.Data[i].Name,
				Enabled:     true,
			}
		}

		err = json.NewEncoder(w).Encode(
			api.BaseResponse[[]*connectorsListElement]{
				Data: &data,
			})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}
