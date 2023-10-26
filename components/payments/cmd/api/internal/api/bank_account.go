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
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type bankAccountsRepository interface {
	ListBankAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.BankAccount, storage.PaginationDetails, error)
}

type bankAccountResponse struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Country       string    `json:"country"`
	Provider      string    `json:"provider"`
	AccountID     string    `json:"accountId,omitempty"`
	Iban          string    `json:"iban,omitempty"`
	AccountNumber string    `json:"accountNumber,omitempty"`
	SwiftBicCode  string    `json:"swiftBicCode,omitempty"`
}

func listBankAccountsHandler(repo bankAccountsRepository) http.HandlerFunc {
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

		ret, paginationDetails, err := repo.ListBankAccounts(r.Context(), pagination)
		if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		data := make([]*bankAccountResponse, len(ret))

		for i := range ret {
			data[i] = &bankAccountResponse{
				ID:        ret[i].ID.String(),
				CreatedAt: ret[i].CreatedAt,
				Country:   ret[i].Country,
				Provider:  ret[i].Provider.String(),
				AccountID: ret[i].AccountID.String(),
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

type readBankAccountRepository interface {
	GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error)
}

func readBankAccountHandler(repo readBankAccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		bankAccountID, err := uuid.Parse(mux.Vars(r)["bankAccountID"])
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		account, err := repo.GetBankAccount(r.Context(), bankAccountID, true)
		if err != nil {
			handleStorageErrors(w, r, err)
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
			Provider:      account.Provider.String(),
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
