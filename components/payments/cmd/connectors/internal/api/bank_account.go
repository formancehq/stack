package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

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

type createBankAccountRepository interface {
	UpsertAccounts(ctx context.Context, provider models.ConnectorProvider, accounts []*models.Account) error
	CreateBankAccount(ctx context.Context, account *models.BankAccount) error
	LinkBankAccountWithAccount(ctx context.Context, id uuid.UUID, accountID *models.AccountID) error
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
			accountID := models.AccountID{
				Reference: bankAccount.ID.String(),
				Provider:  provider,
			}
			err = repo.UpsertAccounts(r.Context(), provider, []*models.Account{
				{
					ID:          accountID,
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

			err = repo.LinkBankAccountWithAccount(r.Context(), bankAccount.ID, &accountID)
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
