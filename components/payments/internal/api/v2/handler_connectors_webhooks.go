package v2

import (
	"io"
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/google/uuid"
)

func connectorsWebhooks(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_connectorsWebhooks")
		defer span.End()

		connectorID, err := models.ConnectorIDFromString(connectorID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		headers := r.Header
		queryValues := r.URL.Query()
		path := r.URL.Path

		webhook := models.Webhook{
			ID:          uuid.New().String(),
			ConnectorID: connectorID,
			QueryValues: queryValues,
			Headers:     headers,
			Body:        body,
		}

		err = backend.ConnectorsHandleWebhooks(ctx, path, webhook)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.RawOk(w, nil)
	}
}
