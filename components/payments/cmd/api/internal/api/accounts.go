package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
)

type accountResponse struct {
	ID              string      `json:"id"`
	Reference       string      `json:"reference"`
	CreatedAt       time.Time   `json:"createdAt"`
	ConnectorID     string      `json:"connectorID"`
	Provider        string      `json:"provider"`
	DefaultCurrency string      `json:"defaultCurrency"` // Deprecated: should be removed soon
	DefaultAsset    string      `json:"defaultAsset"`
	AccountName     string      `json:"accountName"`
	Type            string      `json:"type"`
	Raw             interface{} `json:"raw"`
}

func listAccountsHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		pagination, err := getPagination(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		ret, paginationDetails, err := b.GetService().ListAccounts(r.Context(), pagination)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]*accountResponse, len(ret))

		for i := range ret {
			accountType := ret[i].Type
			if accountType == models.AccountTypeExternalFormance {
				accountType = models.AccountTypeExternal
			}

			data[i] = &accountResponse{
				ID:              ret[i].ID.String(),
				Reference:       ret[i].Reference,
				CreatedAt:       ret[i].CreatedAt,
				ConnectorID:     ret[i].ConnectorID.String(),
				DefaultCurrency: ret[i].DefaultAsset.String(),
				DefaultAsset:    ret[i].DefaultAsset.String(),
				AccountName:     ret[i].AccountName,
				Type:            accountType.String(),
				Raw:             ret[i].RawData,
			}

			if ret[i].Connector != nil {
				data[i].Provider = ret[i].Connector.Provider.String()
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*accountResponse]{
			Cursor: &api.Cursor[*accountResponse]{
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

func readAccountHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		accountID := mux.Vars(r)["accountID"]

		account, err := b.GetService().GetAccount(r.Context(), accountID)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		accountType := account.Type
		if accountType == models.AccountTypeExternalFormance {
			accountType = models.AccountTypeExternal
		}

		data := &accountResponse{
			ID:              account.ID.String(),
			Reference:       account.Reference,
			CreatedAt:       account.CreatedAt,
			ConnectorID:     account.ConnectorID.String(),
			DefaultCurrency: account.DefaultAsset.String(),
			DefaultAsset:    account.DefaultAsset.String(),
			AccountName:     account.AccountName,
			Type:            accountType.String(),
			Raw:             account.RawData,
		}

		if account.Connector != nil {
			data.Provider = account.Connector.Provider.String()
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[accountResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

	}
}
