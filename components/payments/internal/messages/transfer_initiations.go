package messages

import (
	"math/big"
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type transferInitiationsPaymentsMessagePayload struct {
	TransferInitiationID string    `json:"transferInitiationId"`
	PaymentID            string    `json:"paymentId"`
	CreatedAt            time.Time `json:"createdAt"`
	Status               string    `json:"status"`
	Error                string    `json:"error"`
}

type transferInitiationsMessagePayload struct {
	ID                   string                                       `json:"id"`
	CreatedAt            time.Time                                    `json:"createdAt"`
	ScheduleAt           time.Time                                    `json:"scheduledAt"`
	ConnectorID          string                                       `json:"connectorId"`
	Provider             string                                       `json:"provider"`
	Description          string                                       `json:"description"`
	Type                 string                                       `json:"type"`
	SourceAccountID      string                                       `json:"sourceAccountId"`
	DestinationAccountID string                                       `json:"destinationAccountId"`
	Amount               *big.Int                                     `json:"amount"`
	Asset                models.Asset                                 `json:"asset"`
	Attempts             int                                          `json:"attempts"`
	Status               string                                       `json:"status"`
	Error                string                                       `json:"error"`
	RelatedPayments      []*transferInitiationsPaymentsMessagePayload `json:"relatedPayments"`
}

func (m *Messages) NewEventSavedTransferInitiations(tf *models.TransferInitiation) publish.EventMessage {
	payload := transferInitiationsMessagePayload{
		ID:                   tf.ID.String(),
		CreatedAt:            tf.CreatedAt,
		ScheduleAt:           tf.ScheduledAt,
		ConnectorID:          tf.ConnectorID.String(),
		Provider:             tf.Provider.String(),
		Description:          tf.Description,
		Type:                 tf.Type.String(),
		SourceAccountID:      tf.SourceAccountID.String(),
		DestinationAccountID: tf.DestinationAccountID.String(),
		Amount:               tf.Amount,
		Asset:                tf.Asset,
		Attempts:             len(tf.RelatedAdjustments),
	}

	if len(tf.RelatedAdjustments) > 0 {
		// Take the status and error from the last adjustment
		payload.Status = tf.RelatedAdjustments[0].Status.String()
		payload.Error = tf.RelatedAdjustments[0].Error
	}

	payload.RelatedPayments = make([]*transferInitiationsPaymentsMessagePayload, len(tf.RelatedPayments))
	for i, p := range tf.RelatedPayments {
		payload.RelatedPayments[i] = &transferInitiationsPaymentsMessagePayload{
			TransferInitiationID: p.TransferInitiationID.String(),
			PaymentID:            p.PaymentID.String(),
			CreatedAt:            p.CreatedAt,
			Status:               p.Status.String(),
			Error:                p.Error,
		}
	}

	return publish.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedTransferInitiation,
		Payload: payload,
	}
}

type deleteTransferInitiationMessagePayload struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
}

func (m *Messages) NewEventDeleteTransferInitiation(id models.TransferInitiationID) publish.EventMessage {
	return publish.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeDeleteTransferInitiation,
		Payload: deleteTransferInitiationMessagePayload{
			CreatedAt: time.Now().UTC(),
			ID:        id.String(),
		},
	}
}
