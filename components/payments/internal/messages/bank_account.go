package messages

import (
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type bankAccountMessagePayload struct {
	ID            string                         `json:"id"`
	CreatedAt     time.Time                      `json:"createdAt"`
	Name          string                         `json:"name"`
	AccountNumber string                         `json:"accountNumber"`
	IBAN          string                         `json:"iban"`
	SwiftBicCode  string                         `json:"swiftBicCode"`
	Country       string                         `json:"country"`
	Adjustments   []bankAccountAdjustmentPayload `json:"adjustments"`
}

type bankAccountAdjustmentPayload struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	AccountID   string    `json:"accountID"`
	ConnectorID string    `json:"connectorID"`
	Provider    string    `json:"provider"`
}

func (m *Messages) NewEventSavedBankAccounts(bankAccount *models.BankAccount) publish.EventMessage {
	bankAccount.Offuscate()

	payload := bankAccountMessagePayload{
		ID:            bankAccount.ID.String(),
		CreatedAt:     bankAccount.CreatedAt,
		Name:          bankAccount.Name,
		AccountNumber: bankAccount.AccountNumber,
		IBAN:          bankAccount.IBAN,
		SwiftBicCode:  bankAccount.SwiftBicCode,
		Country:       bankAccount.Country,
	}

	for _, adjustment := range bankAccount.Adjustments {
		adjustment := bankAccountAdjustmentPayload{
			ID:          adjustment.ID.String(),
			CreatedAt:   adjustment.CreatedAt,
			AccountID:   adjustment.AccountID.String(),
			Provider:    adjustment.ConnectorID.Provider.String(),
			ConnectorID: adjustment.ConnectorID.String(),
		}

		payload.Adjustments = append(payload.Adjustments, adjustment)
	}

	return publish.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedBankAccount,
		Payload: payload,
	}
}
