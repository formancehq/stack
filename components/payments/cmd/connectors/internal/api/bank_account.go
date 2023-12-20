package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/stack/libs/go-libs/api"
)

type bankAccountResponse struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Country       string    `json:"country"`
	ConnectorID   string    `json:"connectorID"`
	Provider      string    `json:"provider,omitempty"`
	AccountID     string    `json:"accountID,omitempty"`
	Iban          string    `json:"iban,omitempty"`
	AccountNumber string    `json:"accountNumber,omitempty"`
	SwiftBicCode  string    `json:"swiftBicCode,omitempty"`
}

func createBankAccountHandler(
	b backend.ServiceBackend,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var bankAccountRequest service.CreateBankAccountRequest
		err := json.NewDecoder(r.Body).Decode(&bankAccountRequest)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := bankAccountRequest.Validate(); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		bankAccount, err := b.GetService().CreateBankAccount(r.Context(), &bankAccountRequest)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &bankAccountResponse{
			ID:          bankAccount.ID.String(),
			CreatedAt:   bankAccount.CreatedAt,
			Country:     bankAccount.Country,
			ConnectorID: bankAccountRequest.ConnectorID,
			AccountID:   bankAccount.AccountID.String(),
			Provider:    bankAccount.ConnectorID.Provider.String(),
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[bankAccountResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}
