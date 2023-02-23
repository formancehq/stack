package messages

import (
	"time"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/pkg/events"
)

type connectorMessagePayload struct {
	CreatedAt time.Time                `json:"createdAt"`
	Connector models.ConnectorProvider `json:"connector"`
}

func NewEventResetConnector(connector models.ConnectorProvider) events.EventMessage {
	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeConnectorReset,
		Payload: connectorMessagePayload{
			CreatedAt: time.Now().UTC(),
			Connector: connector,
		},
	}
}
