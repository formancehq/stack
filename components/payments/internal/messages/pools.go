package messages

import (
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/google/uuid"
)

type poolMessagePayload struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	AccountIDs []string  `json:"accountIDs"`
}

func NewEventSavedPool(pool *models.Pool) events.EventMessage {
	payload := poolMessagePayload{
		ID:        pool.ID.String(),
		Name:      pool.Name,
		CreatedAt: pool.CreatedAt,
	}

	payload.AccountIDs = make([]string, len(pool.PoolAccounts))
	for i, a := range pool.PoolAccounts {
		payload.AccountIDs[i] = a.AccountID.String()
	}

	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedPool,
		Payload: payload,
	}
}

type deletePoolMessagePayload struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
}

func NewEventDeletePool(id uuid.UUID) events.EventMessage {
	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeDeletePool,
		Payload: deletePoolMessagePayload{
			CreatedAt: time.Now().UTC(),
			ID:        id.String(),
		},
	}
}
