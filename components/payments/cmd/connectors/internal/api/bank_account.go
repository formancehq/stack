package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type bankAccountAdjusmtentsResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	AccountID   string    `json:"accountID"`
	ConnectorID string    `json:"connectorID"`
	Provider    string    `json:"provider"`
}

type bankAccountResponse struct {
	ID            string                            `json:"id"`
	Name          string                            `json:"name"`
	CreatedAt     time.Time                         `json:"createdAt"`
	Country       string                            `json:"country"`
	Iban          string                            `json:"iban,omitempty"`
	AccountNumber string                            `json:"accountNumber,omitempty"`
	SwiftBicCode  string                            `json:"swiftBicCode,omitempty"`
	Metadata      map[string]string                 `json:"metadata,omitempty"`
	Adjustments   []*bankAccountAdjusmtentsResponse `json:"adjustments,omitempty"`

	// Deprecated fields, but clients still use them
	// They correspond to the first adjustment now.
	Provider    string `json:"provider,omitempty"`
	ConnectorID string `json:"connectorID"`
	AccountID   string `json:"accountID,omitempty"`
}

func createBankAccountHandler(
	b backend.ServiceBackend,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "createBankAccountHandler")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")

		var bankAccountRequest service.CreateBankAccountRequest
		err := json.NewDecoder(r.Body).Decode(&bankAccountRequest)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		setAttributesFromRequest(span, &bankAccountRequest)

		if err := bankAccountRequest.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		bankAccount, err := b.GetService().CreateBankAccount(ctx, &bankAccountRequest)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		span.SetAttributes(attribute.String("bankAccount.id", bankAccount.ID.String()))
		span.SetAttributes(attribute.String("bankAccount.createdAt", bankAccount.ID.String()))

		data := &bankAccountResponse{
			ID:        bankAccount.ID.String(),
			Name:      bankAccount.Name,
			CreatedAt: bankAccount.CreatedAt,
			Country:   bankAccount.Country,
			Metadata:  bankAccount.Metadata,
		}

		for _, adjustment := range bankAccount.Adjustments {
			data.Adjustments = append(data.Adjustments, &bankAccountAdjusmtentsResponse{
				ID:          adjustment.ID.String(),
				CreatedAt:   adjustment.CreatedAt,
				AccountID:   adjustment.AccountID.String(),
				ConnectorID: adjustment.ConnectorID.String(),
				Provider:    adjustment.ConnectorID.Provider.String(),
			})
		}

		// Keep compatibility with previous api version
		data.ConnectorID = bankAccountRequest.ConnectorID
		if len(bankAccount.Adjustments) > 0 {
			data.AccountID = bankAccount.Adjustments[0].AccountID.String()
			data.Provider = bankAccount.Adjustments[0].ConnectorID.Provider.String()
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[bankAccountResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func setAttributesFromRequest(span trace.Span, request *service.CreateBankAccountRequest) {
	span.SetAttributes(
		attribute.String("request.name", request.Name),
		attribute.String("request.country", request.Country),
		attribute.String("request.connectorID", request.ConnectorID),
	)
}
