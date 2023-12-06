package messages

import (
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type paymentMessagePayload struct {
	ID          string               `json:"id"`
	Reference   string               `json:"reference"`
	CreatedAt   time.Time            `json:"createdAt"`
	ConnectorID string               `json:"connectorId"`
	Provider    string               `json:"provider"`
	Type        models.PaymentType   `json:"type"`
	Status      models.PaymentStatus `json:"status"`
	Scheme      models.PaymentScheme `json:"scheme"`
	Asset       models.Asset         `json:"asset"`

	// TODO: Remove 'initialAmount' once frontend has switched to 'amount
	InitialAmount *big.Int `json:"initialAmount"`
	Amount        *big.Int `json:"amount"`
}

func NewEventSavedPayments(provider models.ConnectorProvider, payment *models.Payment) events.EventMessage {
	payload := paymentMessagePayload{
		ID:            payment.ID.String(),
		Reference:     payment.Reference,
		Type:          payment.Type,
		Status:        payment.Status,
		InitialAmount: payment.Amount,
		Scheme:        payment.Scheme,
		Asset:         payment.Asset,
		CreatedAt:     payment.CreatedAt,
		Amount:        payment.Amount,
		ConnectorID:   payment.ConnectorID.String(),
		Provider:      provider.String(),
	}

	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedPayments,
		Payload: payload,
	}
}
