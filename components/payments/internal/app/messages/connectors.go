package messages

import (
	"time"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/libs/events/payments"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewEventResetConnector(connector models.ConnectorProvider) (*events.Event, error) {
	now := time.Now().UTC()

	return &events.Event{
		CreatedAt: timestamppb.New(now),
		App:       EventApp,
		Event: &events.Event_ResetConnector{
			ResetConnector: &payments.ResetConnector{
				CreatedAt: timestamppb.New(now),
				Provider:  connector.String(),
			},
		},
	}, nil
}
