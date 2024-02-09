package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type bankAccountRelatedAccountsResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	AccountID   string    `json:"accountID"`
	ConnectorID string    `json:"connectorID"`
	Provider    string    `json:"provider"`
}

type bankAccountResponse struct {
	ID              string                                `json:"id"`
	Name            string                                `json:"name"`
	CreatedAt       time.Time                             `json:"createdAt"`
	Country         string                                `json:"country"`
	Iban            string                                `json:"iban,omitempty"`
	AccountNumber   string                                `json:"accountNumber,omitempty"`
	SwiftBicCode    string                                `json:"swiftBicCode,omitempty"`
	Metadata        map[string]string                     `json:"metadata,omitempty"`
	RelatedAccounts []*bankAccountRelatedAccountsResponse `json:"relatedAccounts,omitempty"`

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

		for _, relatedAccount := range bankAccount.RelatedAccounts {
			data.RelatedAccounts = append(data.RelatedAccounts, &bankAccountRelatedAccountsResponse{
				ID:          relatedAccount.ID.String(),
				CreatedAt:   relatedAccount.CreatedAt,
				AccountID:   relatedAccount.AccountID.String(),
				ConnectorID: relatedAccount.ConnectorID.String(),
				Provider:    relatedAccount.ConnectorID.Provider.String(),
			})
		}

		// Keep compatibility with previous api version
		data.ConnectorID = bankAccountRequest.ConnectorID
		if len(bankAccount.RelatedAccounts) > 0 {
			data.AccountID = bankAccount.RelatedAccounts[0].AccountID.String()
			data.Provider = bankAccount.RelatedAccounts[0].ConnectorID.Provider.String()
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

func forwardBankAccountToConnector(
	b backend.ServiceBackend,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "forwardBankAccountToConnector")
		defer span.End()

		payload := &service.ForwardBankAccountToConnectorRequest{}
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		span.SetAttributes(attribute.String("request.connectorID", payload.ConnectorID))

		if err := payload.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		bankAccountID, ok := mux.Vars(r)["bankAccountID"]
		if !ok {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("bankAccount.id", bankAccountID))

		bankAccount, err := b.GetService().ForwardBankAccountToConnector(ctx, bankAccountID, payload)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := &bankAccountResponse{
			ID:        bankAccount.ID.String(),
			Name:      bankAccount.Name,
			CreatedAt: bankAccount.CreatedAt,
			Country:   bankAccount.Country,
			Metadata:  bankAccount.Metadata,
		}

		for _, relatedAccount := range bankAccount.RelatedAccounts {
			data.RelatedAccounts = append(data.RelatedAccounts, &bankAccountRelatedAccountsResponse{
				ID:          relatedAccount.ID.String(),
				CreatedAt:   relatedAccount.CreatedAt,
				AccountID:   relatedAccount.AccountID.String(),
				ConnectorID: relatedAccount.ConnectorID.String(),
				Provider:    relatedAccount.ConnectorID.Provider.String(),
			})
		}

		// Keep compatibility with previous api version
		data.ConnectorID = payload.ConnectorID
		if len(bankAccount.RelatedAccounts) > 0 {
			data.AccountID = bankAccount.RelatedAccounts[0].AccountID.String()
			data.Provider = bankAccount.RelatedAccounts[0].ConnectorID.Provider.String()
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

func updateBankAccountMetadataHandler(
	b backend.ServiceBackend,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "updateBankAccountMetadataHandler")
		defer span.End()

		payload := &service.UpdateBankAccountMetadataRequest{}
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		for k, v := range payload.Metadata {
			span.SetAttributes(attribute.String(k, v))
		}

		if err := payload.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		bankAccountID, ok := mux.Vars(r)["bankAccountID"]
		if !ok {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("bankAccount.id", bankAccountID))

		err = b.GetService().UpdateBankAccountMetadata(ctx, bankAccountID, payload)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}

func setAttributesFromRequest(span trace.Span, request *service.CreateBankAccountRequest) {
	span.SetAttributes(
		attribute.String("request.name", request.Name),
		attribute.String("request.country", request.Country),
		attribute.String("request.connectorID", request.ConnectorID),
	)
}
