package events

import (
	"time"

	"github.com/formancehq/go-libs/publish"
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

func (e Events) NewEventSavedPool(pool models.Pool) publish.EventMessage {
	payload := poolMessagePayload{
		ID:        pool.ID.String(),
		Name:      pool.Name,
		CreatedAt: pool.CreatedAt,
	}

	payload.AccountIDs = make([]string, len(pool.PoolAccounts))
	for i, a := range pool.PoolAccounts {
		payload.AccountIDs[i] = a.AccountID.String()
	}

	return publish.EventMessage{
		IdempotencyKey: pool.IdempotencyKey(),
		Date:           time.Now().UTC(),
		App:            events.EventApp,
		Version:        events.EventVersion,
		Type:           events.EventTypeSavedPool,
		Payload:        payload,
	}
}

type deletePoolMessagePayload struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
}

func (e Events) NewEventDeletePool(id uuid.UUID) publish.EventMessage {
	return publish.EventMessage{
		IdempotencyKey: id.String(),
		Date:           time.Now().UTC(),
		App:            events.EventApp,
		Version:        events.EventVersion,
		Type:           events.EventTypeDeletePool,
		Payload: deletePoolMessagePayload{
			CreatedAt: time.Now().UTC(),
			ID:        id.String(),
		},
	}
}
