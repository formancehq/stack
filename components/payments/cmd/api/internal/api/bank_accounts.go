package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	// They correspond to the first bank account adjustment now.
	ConnectorID string `json:"connectorID"`
	Provider    string `json:"provider,omitempty"`
	AccountID   string `json:"accountID,omitempty"`
}

func listBankAccountsHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		query, err := bunpaginate.Extract[storage.ListBankAccountQuery](r, func() (*storage.ListBankAccountQuery, error) {
			options, err := getPagination(r, storage.BankAccountQuery{})
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListBankAccountQuery(*options)), nil
		})
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := b.GetService().ListBankAccounts(r.Context(), *query)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		ret := cursor.Data
		data := make([]*bankAccountResponse, len(ret))

		for i := range ret {
			data[i] = &bankAccountResponse{
				ID:        ret[i].ID.String(),
				Name:      ret[i].Name,
				CreatedAt: ret[i].CreatedAt,
				Country:   ret[i].Country,
				Metadata:  ret[i].Metadata,
			}

			// Deprecated fields, but clients still use them
			if len(ret[i].RelatedAccounts) > 0 {
				data[i].ConnectorID = ret[i].RelatedAccounts[0].ConnectorID.String()
				data[i].AccountID = ret[i].RelatedAccounts[0].AccountID.String()
				data[i].Provider = ret[i].RelatedAccounts[0].ConnectorID.Provider.String()
			}

			for _, adjustment := range ret[i].RelatedAccounts {
				data[i].RelatedAccounts = append(data[i].RelatedAccounts, &bankAccountRelatedAccountsResponse{
					ID:          adjustment.ID.String(),
					CreatedAt:   adjustment.CreatedAt,
					AccountID:   adjustment.AccountID.String(),
					ConnectorID: adjustment.ConnectorID.String(),
					Provider:    adjustment.ConnectorID.Provider.String(),
				})
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*bankAccountResponse]{
			Cursor: &api.Cursor[*bankAccountResponse]{
				PageSize: cursor.PageSize,
				HasMore:  cursor.HasMore,
				Previous: cursor.Previous,
				Next:     cursor.Next,
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
			Name:          account.Name,
			CreatedAt:     account.CreatedAt,
			Country:       account.Country,
			Iban:          account.IBAN,
			AccountNumber: account.AccountNumber,
			SwiftBicCode:  account.SwiftBicCode,
			Metadata:      account.Metadata,
		}

		// Deprecated fields, but clients still use them
		if len(account.RelatedAccounts) > 0 {
			data.ConnectorID = account.RelatedAccounts[0].ConnectorID.String()
			data.AccountID = account.RelatedAccounts[0].AccountID.String()
			data.Provider = account.RelatedAccounts[0].ConnectorID.Provider.String()
		}

		for _, adjustment := range account.RelatedAccounts {
			data.RelatedAccounts = append(data.RelatedAccounts, &bankAccountRelatedAccountsResponse{
				ID:          adjustment.ID.String(),
				CreatedAt:   adjustment.CreatedAt,
				AccountID:   adjustment.AccountID.String(),
				ConnectorID: adjustment.ConnectorID.String(),
				Provider:    adjustment.ConnectorID.Provider.String(),
			})
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
