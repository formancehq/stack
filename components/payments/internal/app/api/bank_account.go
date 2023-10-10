package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/storage"
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
						handleValidationError(w, r, errors.New("sort order not well specified, got "+parts[1]))

						return
					}
				}

				column := parts[0]

				sorter.Add(column, order)
			}
		}

		pageSize, err := pageSizeQueryParam(r)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		pagination, err := storage.Paginate(pageSize, r.URL.Query().Get("cursor"), sorter, nil)
		if err != nil {
			handleValidationError(w, r, err)

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
			handleServerError(w, r, err)

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
			handleErrorBadRequest(w, r, err)

			return
		}

		account, err := repo.GetBankAccount(r.Context(), bankAccountID, true)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		if err := account.Offuscate(); err != nil {
			handleServerError(w, r, err)

			return
		}

		data := &bankAccountResponse{
			ID:            account.ID.String(),
			CreatedAt:     account.CreatedAt,
			Country:       account.Country,
			Provider:      account.Provider.String(),
			Iban:          account.IBAN,
			AccountNumber: account.AccountNumber,
			SwiftBicCode:  account.SwiftBicCode,
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[bankAccountResponse]{
			Data: data,
		})
		if err != nil {
			handleServerError(w, r, err)

			return
		}

	}
}

type createBankAccountRepository interface {
	UpsertAccounts(ctx context.Context, provider models.ConnectorProvider, accounts []*models.Account) error
	CreateBankAccount(ctx context.Context, account *models.BankAccount) error
	IsInstalled(ctx context.Context, provider models.ConnectorProvider) (bool, error)
}

type createBankAccountRequest struct {
	AccountNumber string `json:"accountNumber"`
	IBAN          string `json:"iban"`
	SwiftBicCode  string `json:"swiftBicCode"`
	Country       string `json:"country"`
	Provider      string `json:"provider"`
	Name          string `json:"name"`
}

func createBankAccountHandler(repo createBankAccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var bankAccountRequest createBankAccountRequest

		err := json.NewDecoder(r.Body).Decode(&bankAccountRequest)
		if err != nil {
			handleErrorBadRequest(w, r, err)

			return
		}

		if bankAccountRequest.AccountNumber == "" &&
			bankAccountRequest.IBAN == "" {
			handleErrorBadRequest(w, r, errors.New("either accountNumber or iban must be provided"))

			return
		}

		if bankAccountRequest.Name == "" {
			handleErrorBadRequest(w, r, errors.New("name must be provided"))

			return
		}

		provider, err := models.ConnectorProviderFromString(bankAccountRequest.Provider)
		if err != nil {
			handleErrorBadRequest(w, r, err)

			return
		}

		if provider != models.ConnectorProviderBankingCircle {
			// For now, bank accounts can only be created for BankingCircle
			// in the future, we will support other providers
			handleErrorBadRequest(w, r, errors.New("provider not supported"))

			return
		}

		isInstalled, err := repo.IsInstalled(r.Context(), provider)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		if !isInstalled {
			handleErrorBadRequest(w, r, errors.New("connector not installed"))

			return
		}

		bankAccount := &models.BankAccount{
			CreatedAt:     time.Now().UTC(),
			AccountNumber: bankAccountRequest.AccountNumber,
			IBAN:          bankAccountRequest.IBAN,
			SwiftBicCode:  bankAccountRequest.SwiftBicCode,
			Country:       bankAccountRequest.Country,
			Provider:      provider,
			Name:          bankAccountRequest.Name,
		}
		err = repo.CreateBankAccount(r.Context(), bankAccount)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		// BankingCircle does not have external accounts so we need to create
		// one by hand
		if provider == models.ConnectorProviderBankingCircle {
			err = repo.UpsertAccounts(r.Context(), provider, []*models.Account{
				{
					ID: models.AccountID{
						Reference: bankAccount.ID.String(),
						Provider:  provider,
					},
					CreatedAt:   time.Now(),
					Reference:   bankAccount.ID.String(),
					Provider:    provider,
					AccountName: bankAccount.Name,
					Type:        models.AccountTypeExternalFormance,
				},
			})
			if err != nil {
				handleStorageErrors(w, r, err)

				return
			}
		}

		data := &bankAccountResponse{
			ID:        bankAccount.ID.String(),
			CreatedAt: bankAccount.CreatedAt,
			Country:   bankAccount.Country,
			Provider:  bankAccount.Provider.String(),
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[bankAccountResponse]{
			Data: data,
		})
		if err != nil {
			handleServerError(w, r, err)

			return
		}
	}
}
