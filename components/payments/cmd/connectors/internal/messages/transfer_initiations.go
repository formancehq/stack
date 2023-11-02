package messages

import (
	"math/big"
	"time"

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
	UpdatedAt            time.Time                                    `json:"updatedAt"`
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

func NewEventSavedTransferInitiations(tf *models.TransferInitiation) events.EventMessage {
	payload := transferInitiationsMessagePayload{
		ID:                   tf.ID.String(),
		CreatedAt:            tf.CreatedAt,
		ScheduleAt:           tf.ScheduledAt,
		UpdatedAt:            tf.UpdatedAt,
		Provider:             tf.Provider.String(),
		Description:          tf.Description,
		Type:                 tf.Type.String(),
		SourceAccountID:      tf.SourceAccountID.String(),
		DestinationAccountID: tf.DestinationAccountID.String(),
		Amount:               tf.Amount,
		Asset:                tf.Asset,
		Attempts:             tf.Attempts,
		Status:               tf.Status.String(),
		Error:                tf.Error,
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

	return events.EventMessage{
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

func NewEventDeleteTransferInitiation(id models.TransferInitiationID) events.EventMessage {
	return events.EventMessage{
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
