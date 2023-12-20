package messages

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type link struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

type paymentMessagePayload struct {
	ID                   string               `json:"id"`
	Reference            string               `json:"reference"`
	CreatedAt            time.Time            `json:"createdAt"`
	ConnectorID          string               `json:"connectorId"`
	Provider             string               `json:"provider"`
	Type                 models.PaymentType   `json:"type"`
	Status               models.PaymentStatus `json:"status"`
	Scheme               models.PaymentScheme `json:"scheme"`
	Asset                models.Asset         `json:"asset"`
	SourceAccountID      *models.AccountID    `json:"sourceAccountId,omitempty"`
	DestinationAccountID *models.AccountID    `json:"destinationAccountId,omitempty"`
	Links                []link               `json:"links"`
	RawData              json.RawMessage      `json:"rawData"`

	// TODO: Remove 'initialAmount' once frontend has switched to 'amount
	InitialAmount *big.Int          `json:"initialAmount"`
	Amount        *big.Int          `json:"amount"`
	Metadata      map[string]string `json:"metadata"`
}

func (m *Messages) NewEventSavedPayments(provider models.ConnectorProvider, payment *models.Payment) publish.EventMessage {
	payload := paymentMessagePayload{
		ID:                   payment.ID.String(),
		Reference:            payment.Reference,
		Type:                 payment.Type,
		Status:               payment.Status,
		InitialAmount:        payment.InitialAmount,
		Amount:               payment.Amount,
		Scheme:               payment.Scheme,
		Asset:                payment.Asset,
		CreatedAt:            payment.CreatedAt,
		ConnectorID:          payment.ConnectorID.String(),
		Provider:             provider.String(),
		SourceAccountID:      payment.SourceAccountID,
		DestinationAccountID: payment.DestinationAccountID,

		RawData: payment.RawData,
		Metadata: func() map[string]string {
			ret := make(map[string]string)
			for _, m := range payment.Metadata {
				ret[m.Key] = m.Value
			}
			return ret
		}(),
	}

	if payment.SourceAccountID != nil {
		payload.Links = append(payload.Links, link{
			Name: "source_account",
			URI:  m.stackURL + "/api/payments/accounts/" + payment.SourceAccountID.String(),
		})
	}

	if payment.DestinationAccountID != nil {
		payload.Links = append(payload.Links, link{
			Name: "destination_account",
			URI:  m.stackURL + "/api/payments/accounts/" + payment.DestinationAccountID.String(),
		})
	}

	return publish.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedPayments,
		Payload: payload,
	}
}
