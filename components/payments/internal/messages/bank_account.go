package messages

import (
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type bankAccountMessagePayload struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	ConnectorID   string    `json:"connectorId"`
	Provider      string    `json:"provider"`
	Name          string    `json:"name"`
	AccountNumber string    `json:"accountNumber"`
	IBAN          string    `json:"iban"`
	SwiftBicCode  string    `json:"swiftBicCode"`
	Country       string    `json:"country"`
	AccountID     string    `json:"accountId"`
}

func (m *Messages) NewEventSavedBankAccounts(bankAccount *models.BankAccount) publish.EventMessage {
	bankAccount.Offuscate()

	payload := bankAccountMessagePayload{
		ID:            bankAccount.ID.String(),
		CreatedAt:     bankAccount.CreatedAt,
		ConnectorID:   bankAccount.ConnectorID.String(),
		Provider:      bankAccount.ConnectorID.Provider.String(),
		Name:          bankAccount.Name,
		AccountNumber: bankAccount.AccountNumber,
		IBAN:          bankAccount.IBAN,
		SwiftBicCode:  bankAccount.SwiftBicCode,
		Country:       bankAccount.Country,
		AccountID:     bankAccount.AccountID.String(),
	}

	return publish.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedBankAccount,
		Payload: payload,
	}
}
