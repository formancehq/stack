package api

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type readConnectorsResponseElement struct {
	Provider    models.ConnectorProvider `json:"provider" bson:"provider"`
	ConnectorID string                   `json:"connectorID" bson:"connectorID"`
	Name        string                   `json:"name" bson:"name"`
	Enabled     bool                     `json:"enabled" bson:"enabled"`
}

func readConnectorsHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "readConnectorsHandler")
		defer span.End()

		res, err := b.GetService().ListConnectors(ctx)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		span.SetAttributes(attribute.Int("count", len(res)))

		data := make([]readConnectorsResponseElement, len(res))

		for i := range res {
			data[i] = readConnectorsResponseElement{
				Provider:    res[i].Provider,
				ConnectorID: res[i].ID.String(),
				Name:        res[i].Name,
				Enabled:     true,
			}
		}

		err = json.NewEncoder(w).Encode(
			api.BaseResponse[[]readConnectorsResponseElement]{
				Data: &data,
			})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}
