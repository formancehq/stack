package events

import (
	"time"

	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type connectorMessagePayload struct {
	CreatedAt   time.Time `json:"createdAt"`
	ConnectorID string    `json:"connectorId"`
}

func (e Events) NewEventResetConnector(connectorID models.ConnectorID) publish.EventMessage {
	return publish.EventMessage{
		IdempotencyKey: connectorID.String(),
		Date:           time.Now().UTC(),
		App:            events.EventApp,
		Version:        events.EventVersion,
		Type:           events.EventTypeConnectorReset,
		Payload: connectorMessagePayload{
			CreatedAt:   time.Now().UTC(),
			ConnectorID: connectorID.String(),
		},
	}
}
