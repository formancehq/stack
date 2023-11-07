package messages

import (
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type bankAccountMessagePayload struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	ConnectorID   string    `json:"connectorId"`
	Name          string    `json:"name"`
	AccountNumber string    `json:"accountNumber"`
	IBAN          string    `json:"iban"`
	SwiftBicCode  string    `json:"swiftBicCode"`
	Country       string    `json:"country"`
	AccountID     string    `json:"accountId"`
}

func NewEventSavedBankAccounts(bankAccount *models.BankAccount) events.EventMessage {
	bankAccount.Offuscate()

	payload := bankAccountMessagePayload{
		ID:            bankAccount.ID.String(),
		CreatedAt:     bankAccount.CreatedAt,
		ConnectorID:   bankAccount.ConnectorID.String(),
		Name:          bankAccount.Name,
		AccountNumber: bankAccount.AccountNumber,
		IBAN:          bankAccount.IBAN,
		SwiftBicCode:  bankAccount.SwiftBicCode,
		Country:       bankAccount.Country,
		AccountID:     bankAccount.AccountID.String(),
	}

	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedBankAccount,
		Payload: payload,
	}
}
