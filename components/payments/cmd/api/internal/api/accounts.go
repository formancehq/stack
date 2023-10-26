package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type accountsRepository interface {
	ListAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Account, storage.PaginationDetails, error)
}

type accountResponse struct {
	ID              string      `json:"id"`
	Reference       string      `json:"reference"`
	CreatedAt       time.Time   `json:"createdAt"`
	Provider        string      `json:"provider"`
	DefaultCurrency string      `json:"defaultCurrency"` // Deprecated: should be removed soon
	DefaultAsset    string      `json:"defaultAsset"`
	AccountName     string      `json:"accountName"`
	Type            string      `json:"type"`
	Raw             interface{} `json:"raw"`
}

func listAccountsHandler(repo accountsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var sorter storage.Sorter

		if sortParams := r.URL.Query()["sort"]; sortParams != nil {
			for _, s := range sortParams {
				parts := strings.SplitN(s, ":", 2)

				var order storage.SortOrder

				if len(parts) > 1 {
					//nolint:goconst // allow duplicate string
					switch parts[1] {
					case "asc", "ASC":
						order = storage.SortOrderAsc
					case "dsc", "desc", "DSC", "DESC":
						order = storage.SortOrderDesc
					default:
						api.BadRequest(w, ErrValidation, errors.New("sort order not well specified, got "+parts[1]))

						return
					}
				}

				column := parts[0]

				sorter.Add(column, order)
			}
		}

		pageSize, err := pageSizeQueryParam(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		pagination, err := storage.Paginate(pageSize, r.URL.Query().Get("cursor"), sorter, nil)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		ret, paginationDetails, err := repo.ListAccounts(r.Context(), pagination)
		if err != nil {
			handleStorageErrors(w, r, err)
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
				Provider:        ret[i].Provider.String(),
				DefaultCurrency: ret[i].DefaultAsset.String(),
				DefaultAsset:    ret[i].DefaultAsset.String(),
				AccountName:     ret[i].AccountName,
				Type:            accountType.String(),
				Raw:             ret[i].RawData,
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

type readAccountRepository interface {
	GetAccount(ctx context.Context, id string) (*models.Account, error)
}

func readAccountHandler(repo readAccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		accountID := mux.Vars(r)["accountID"]

		account, err := repo.GetAccount(r.Context(), accountID)
		if err != nil {
			handleStorageErrors(w, r, err)

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
			Provider:        account.Provider.String(),
			DefaultCurrency: account.DefaultAsset.String(),
			DefaultAsset:    account.DefaultAsset.String(),
			AccountName:     account.AccountName,
			Type:            accountType.String(),
			Raw:             account.RawData,
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
