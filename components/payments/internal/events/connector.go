package events

import (
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

type connectorMessagePayload struct {
	CreatedAt   time.Time `json:"createdAt"`
	ConnectorID string    `json:"connectorId"`
}

func (e Events) NewEventResetConnector(connectorID models.ConnectorID) publish.EventMessage {
	return publish.EventMessage{
		IdempotemcyKey: connectorID.String(),
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
