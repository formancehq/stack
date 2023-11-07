package messages

import (
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type accountMessagePayload struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Reference    string    `json:"reference"`
	ConnectorID  string    `json:"connectorId"`
	Provider     string    `json:"provider"`
	DefaultAsset string    `json:"defaultAsset"`
	AccountName  string    `json:"accountName"`
	Type         string    `json:"type"`
}

func NewEventSavedAccounts(provider models.ConnectorProvider, account *models.Account) events.EventMessage {
	payload := accountMessagePayload{
		ID:           account.ID.String(),
		CreatedAt:    account.CreatedAt,
		Reference:    account.Reference,
		ConnectorID:  account.ConnectorID.String(),
		DefaultAsset: account.DefaultAsset.String(),
		AccountName:  account.AccountName,
		Type:         string(account.Type),
		Provider:     provider.String(),
	}

	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedAccounts,
		Payload: payload,
	}
}
