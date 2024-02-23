package api

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"

	"github.com/gorilla/mux"
)

func updateMetadataHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "updateMetadataHandler")
		defer span.End()

		paymentID, err := models.PaymentIDFromString(mux.Vars(r)["paymentID"])
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("request.paymentID", paymentID.String()))

		var metadata service.UpdateMetadataRequest
		if r.ContentLength == 0 {
			var err = errors.New("body is required")
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&metadata)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		for k, v := range metadata {
			span.SetAttributes(attribute.String("request.metadata."+k, v))
		}

		err = b.GetService().UpdatePaymentMetadata(ctx, *paymentID, metadata)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
