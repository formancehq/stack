package events

import (
	"time"

	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type bankAccountMessagePayload struct {
	ID              string                              `json:"id"`
	CreatedAt       time.Time                           `json:"createdAt"`
	Name            string                              `json:"name"`
	AccountNumber   string                              `json:"accountNumber"`
	IBAN            string                              `json:"iban"`
	SwiftBicCode    string                              `json:"swiftBicCode"`
	Country         string                              `json:"country"`
	RelatedAccounts []bankAccountRelatedAccountsPayload `json:"adjustments"`
}

type bankAccountRelatedAccountsPayload struct {
	CreatedAt   time.Time `json:"createdAt"`
	AccountID   string    `json:"accountID"`
	ConnectorID string    `json:"connectorID"`
	Provider    string    `json:"provider"`
}

func (e Events) NewEventSavedBankAccounts(bankAccount models.BankAccount) publish.EventMessage {
	bankAccount.Offuscate()

	payload := bankAccountMessagePayload{
		ID:        bankAccount.ID.String(),
		CreatedAt: bankAccount.CreatedAt,
		Name:      bankAccount.Name,
	}

	if bankAccount.AccountNumber != nil {
		payload.AccountNumber = *bankAccount.AccountNumber
	}

	if bankAccount.IBAN != nil {
		payload.IBAN = *bankAccount.IBAN
	}

	if bankAccount.SwiftBicCode != nil {
		payload.SwiftBicCode = *bankAccount.SwiftBicCode
	}

	if bankAccount.Country != nil {
		payload.Country = *bankAccount.Country
	}

	for _, relatedAccount := range bankAccount.RelatedAccounts {
		relatedAccount := bankAccountRelatedAccountsPayload{
			CreatedAt:   relatedAccount.CreatedAt,
			AccountID:   relatedAccount.AccountID.String(),
			Provider:    relatedAccount.ConnectorID.Provider,
			ConnectorID: relatedAccount.ConnectorID.String(),
		}

		payload.RelatedAccounts = append(payload.RelatedAccounts, relatedAccount)
	}

	return publish.EventMessage{
		IdempotencyKey: bankAccount.IdempotencyKey(),
		Date:           time.Now().UTC(),
		App:            events.EventApp,
		Version:        events.EventVersion,
		Type:           events.EventTypeSavedBankAccount,
		Payload:        payload,
	}
}
