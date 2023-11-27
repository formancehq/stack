package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type bankAccountResponse struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Country       string    `json:"country"`
	ConnectorID   string    `json:"connectorID"`
	AccountID     string    `json:"accountId,omitempty"`
	Iban          string    `json:"iban,omitempty"`
	AccountNumber string    `json:"accountNumber,omitempty"`
	SwiftBicCode  string    `json:"swiftBicCode,omitempty"`
}

func listBankAccountsHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		pagination, err := getPagination(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		ret, paginationDetails, err := b.GetService().ListBankAccounts(r.Context(), pagination)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]*bankAccountResponse, len(ret))

		for i := range ret {
			data[i] = &bankAccountResponse{
				ID:          ret[i].ID.String(),
				CreatedAt:   ret[i].CreatedAt,
				Country:     ret[i].Country,
				ConnectorID: ret[i].ConnectorID.String(),
				AccountID:   ret[i].AccountID.String(),
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*bankAccountResponse]{
			Cursor: &api.Cursor[*bankAccountResponse]{
				PageSize: paginationDetails.PageSize,
				HasMore:  paginationDetails.HasMore,
				Previous: paginationDetails.PreviousPage,
				Next:     paginationDetails.NextPage,
				Data:     data,
			},
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func readBankAccountHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		bankAccountID, err := uuid.Parse(mux.Vars(r)["bankAccountID"])
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		account, err := b.GetService().GetBankAccount(r.Context(), bankAccountID, true)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		if err := account.Offuscate(); err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		data := &bankAccountResponse{
			ID:            account.ID.String(),
			CreatedAt:     account.CreatedAt,
			Country:       account.Country,
			ConnectorID:   account.ConnectorID.String(),
			AccountID:     account.AccountID.String(),
			Iban:          account.IBAN,
			AccountNumber: account.AccountNumber,
			SwiftBicCode:  account.SwiftBicCode,
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
