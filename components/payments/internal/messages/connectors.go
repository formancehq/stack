package messages

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

func (m *Messages) NewEventResetConnector(connectorID models.ConnectorID) publish.EventMessage {
	return publish.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeConnectorReset,
		Payload: connectorMessagePayload{
			CreatedAt:   time.Now().UTC(),
			ConnectorID: connectorID.String(),
		},
	}
}
