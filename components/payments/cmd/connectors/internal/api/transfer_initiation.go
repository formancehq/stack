package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type transferInitiationResponse struct {
	ID                   string            `json:"id"`
	CreatedAt            time.Time         `json:"createdAt"`
	ScheduledAt          time.Time         `json:"scheduledAt"`
	Description          string            `json:"description"`
	SourceAccountID      string            `json:"sourceAccountID"`
	DestinationAccountID string            `json:"destinationAccountID"`
	ConnectorID          string            `json:"connectorID"`
	Type                 string            `json:"type"`
	Amount               *big.Int          `json:"amount"`
	Asset                string            `json:"asset"`
	Status               string            `json:"status"`
	Error                string            `json:"error"`
	Metadata             map[string]string `json:"metadata"`
}

func createTransferInitiationHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "createTransferInitiationHandler")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")

		payload := &service.CreateTransferInitiationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		setSpanAttributesFromRequest(span, payload)

		if err := payload.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		tf, err := b.GetService().CreateTransferInitiation(ctx, payload)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		span.SetAttributes(
			attribute.String("transfer.id", tf.ID.String()),
			attribute.String("transfer.createdAt", tf.CreatedAt.String()),
			attribute.String("connectorID", tf.ConnectorID.String()),
		)

		data := &transferInitiationResponse{
			ID:                   tf.ID.String(),
			CreatedAt:            tf.CreatedAt,
			ScheduledAt:          tf.ScheduledAt,
			Description:          tf.Description,
			SourceAccountID:      tf.SourceAccountID.String(),
			DestinationAccountID: tf.DestinationAccountID.String(),
			ConnectorID:          tf.ConnectorID.String(),
			Type:                 tf.Type.String(),
			Amount:               tf.Amount,
			Asset:                tf.Asset.String(),
			Metadata:             tf.Metadata,
		}

		if len(tf.RelatedAdjustments) > 0 {
			// Take the status and error from the last adjustment
			data.Status = tf.RelatedAdjustments[0].Status.String()
			data.Error = tf.RelatedAdjustments[0].Error
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[transferInitiationResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func updateTransferInitiationStatusHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "updateTransferInitiationStatusHandler")
		defer span.End()

		payload := &service.UpdateTransferInitiationStatusRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		span.SetAttributes(attribute.String("request.status", payload.Status))

		if err := payload.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		transferID, ok := mux.Vars(r)["transferID"]
		if !ok {
			otel.RecordError(span, errors.New("missing transferID"))
			api.BadRequest(w, ErrInvalidID, errors.New("missing transferID"))
			return
		}

		span.SetAttributes(attribute.String("transfer.id", transferID))

		if err := b.GetService().UpdateTransferInitiationStatus(ctx, transferID, payload); err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func retryTransferInitiationHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "retryTransferInitiationHandler")
		defer span.End()

		transferID, ok := mux.Vars(r)["transferID"]
		if !ok {
			otel.RecordError(span, errors.New("missing transferID"))
			api.BadRequest(w, ErrInvalidID, errors.New("missing transferID"))
			return
		}

		span.SetAttributes(attribute.String("transfer.id", transferID))

		if err := b.GetService().RetryTransferInitiation(ctx, transferID); err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteTransferInitiationHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "deleteTransferInitiationHandler")
		defer span.End()

		transferID, ok := mux.Vars(r)["transferID"]
		if !ok {
			otel.RecordError(span, errors.New("missing transferID"))
			api.BadRequest(w, ErrInvalidID, errors.New("missing transferID"))
			return
		}

		span.SetAttributes(attribute.String("transfer.id", transferID))

		if err := b.GetService().DeleteTransferInitiation(ctx, transferID); err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func setSpanAttributesFromRequest(span trace.Span, transfer *service.CreateTransferInitiationRequest) {
	span.SetAttributes(
		attribute.String("request.reference", transfer.Reference),
		attribute.String("request.scheduledAt", transfer.ScheduledAt.String()),
		attribute.String("request.description", transfer.Description),
		attribute.String("request.sourceAccountID", transfer.SourceAccountID),
		attribute.String("request.destinationAccountID", transfer.DestinationAccountID),
		attribute.String("request.connectorID", transfer.ConnectorID),
		attribute.String("request.provider", transfer.Provider),
		attribute.String("request.type", transfer.Type),
		attribute.String("request.amount", transfer.Amount.String()),
		attribute.String("request.asset", transfer.Asset),
		attribute.String("request.validated", transfer.Asset),
	)
}
