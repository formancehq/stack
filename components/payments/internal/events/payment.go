package events

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type paymentMessagePayload struct {
	ID                   string          `json:"id"`
	ConnectorID          string          `json:"connectorId"`
	Provider             string          `json:"provider"`
	Reference            string          `json:"reference"`
	CreatedAt            time.Time       `json:"createdAt"`
	Type                 string          `json:"type"`
	Status               string          `json:"status"`
	Scheme               string          `json:"scheme"`
	Asset                string          `json:"asset"`
	SourceAccountID      string          `json:"sourceAccountId,omitempty"`
	DestinationAccountID string          `json:"destinationAccountId,omitempty"`
	Links                []api.Link      `json:"links"`
	RawData              json.RawMessage `json:"rawData"`

	Amount   *big.Int          `json:"amount"`
	Metadata map[string]string `json:"metadata"`
}

func (e Events) NewEventSavedPayments(payment models.Payment, adjustment models.PaymentAdjustment) publish.EventMessage {
	payload := paymentMessagePayload{
		ID:          payment.ID.String(),
		Reference:   payment.Reference,
		Type:        payment.Type.String(),
		Status:      payment.Status.String(),
		Amount:      payment.Amount,
		Scheme:      payment.Scheme.String(),
		Asset:       payment.Asset,
		CreatedAt:   payment.CreatedAt,
		ConnectorID: payment.ConnectorID.String(),
		Provider:    payment.ConnectorID.Provider,
		SourceAccountID: func() string {
			if payment.SourceAccountID == nil {
				return ""
			}
			return payment.SourceAccountID.String()
		}(),
		DestinationAccountID: func() string {
			if payment.DestinationAccountID == nil {
				return ""
			}
			return payment.DestinationAccountID.String()
		}(),
		RawData:  adjustment.Raw,
		Metadata: payment.Metadata,
	}

	if payment.SourceAccountID != nil {
		payload.Links = append(payload.Links, api.Link{
			Name: "source_account",
			URI:  e.stackURL + "/api/payments/accounts/" + payment.SourceAccountID.String(),
		})
	}

	if payment.DestinationAccountID != nil {
		payload.Links = append(payload.Links, api.Link{
			Name: "destination_account",
			URI:  e.stackURL + "/api/payments/accounts/" + payment.DestinationAccountID.String(),
		})
	}

	return publish.EventMessage{
		IdempotencyKey: adjustment.IdempotencyKey(),
		Date:           time.Now().UTC(),
		App:            events.EventApp,
		Version:        events.EventVersion,
		Type:           events.EventTypeSavedPayments,
		Payload:        payload,
	}
}
