package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func connectorsConfigs(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.Tracer().Start(r.Context(), "v2_connectorsConfigs")
		defer span.End()

		configs := backend.ConnectorsConfigs()

		err := json.NewEncoder(w).Encode(api.BaseResponse[plugins.Configs]{
			Data: &configs,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}
