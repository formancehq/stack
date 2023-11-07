package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/connectors/internal/messages"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

type createBankAccountRepository interface {
	UpsertAccounts(ctx context.Context, accounts []*models.Account) error
	CreateBankAccount(ctx context.Context, account *models.BankAccount) error
	LinkBankAccountWithAccount(ctx context.Context, id uuid.UUID, accountID *models.AccountID) error
	GetConnector(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error)
}

type createBankAccountRequest struct {
	AccountNumber string `json:"accountNumber"`
	IBAN          string `json:"iban"`
	SwiftBicCode  string `json:"swiftBicCode"`
	Country       string `json:"country"`
	ConnectorID   string `json:"connectorID"`
	Name          string `json:"name"`
}

func createBankAccountHandler(
	repo createBankAccountRepository,
	publisher message.Publisher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var bankAccountRequest createBankAccountRequest

		err := json.NewDecoder(r.Body).Decode(&bankAccountRequest)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if bankAccountRequest.AccountNumber == "" &&
			bankAccountRequest.IBAN == "" {
			api.BadRequest(w, ErrValidation, errors.New("either accountNumber or iban must be provided"))
			return
		}

		if bankAccountRequest.Name == "" {
			api.BadRequest(w, ErrValidation, errors.New("name must be provided"))
			return
		}

		connectorID, err := models.ConnectorIDFromString(bankAccountRequest.ConnectorID)
		if err != nil {
			api.BadRequest(w, ErrValidation, errors.New("connectorId must be a valid"))
			return
		}

		connector, err := repo.GetConnector(r.Context(), connectorID)
		if err != nil && errors.Is(err, storage.ErrNotFound) {
			api.BadRequest(w, ErrValidation, errors.New("connector not installed"))
			return
		} else if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		if connector.Provider != models.ConnectorProviderBankingCircle {
			// For now, bank accounts can only be created for BankingCircle
			// in the future, we will support other providers
			api.BadRequest(w, ErrValidation, errors.New("provider not supported"))
			return
		}

		bankAccount := &models.BankAccount{
			CreatedAt:     time.Now().UTC(),
			AccountNumber: bankAccountRequest.AccountNumber,
			IBAN:          bankAccountRequest.IBAN,
			SwiftBicCode:  bankAccountRequest.SwiftBicCode,
			Country:       bankAccountRequest.Country,
			ConnectorID:   connectorID,
			Name:          bankAccountRequest.Name,
		}
		err = repo.CreateBankAccount(r.Context(), bankAccount)
		if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		// BankingCircle does not have external accounts so we need to create
		// one by hand
		if connector.Provider == models.ConnectorProviderBankingCircle {
			accountID := models.AccountID{
				Reference:   bankAccount.ID.String(),
				ConnectorID: connector.ID,
			}
			err = repo.UpsertAccounts(r.Context(), []*models.Account{
				{
					ID:          accountID,
					CreatedAt:   time.Now(),
					Reference:   bankAccount.ID.String(),
					ConnectorID: connector.ID,
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

			bankAccount.AccountID = &accountID
		}

		if err := publisher.Publish(
			events.TopicPayments,
			publish.NewMessage(
				r.Context(),
				messages.NewEventSavedBankAccounts(bankAccount),
			),
		); err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		data := &bankAccountResponse{
			ID:          bankAccount.ID.String(),
			CreatedAt:   bankAccount.CreatedAt,
			Country:     bankAccount.Country,
			ConnectorID: connectorID.String(),
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
