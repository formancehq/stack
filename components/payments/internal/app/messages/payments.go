package messages

import (
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/libs/events/payments"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	TopicPayments   = "payments"
	TopicConnectors = "connectors"

	EventApp = "payments"
)

func NewEventSavedPayments(payment *models.Payment, provider models.ConnectorProvider) (*events.Event, error) {
	return &events.Event{
		CreatedAt: timestamppb.New(payment.CreatedAt),
		App:       EventApp,
		Event: &events.Event_PaymentSaved{
			PaymentSaved: &payments.PaymentSaved{
				Id:        payment.ID.String(),
				Reference: payment.Reference,
				CreatedAt: timestamppb.New(payment.CreatedAt),
				Provider:  provider.String(),
				Type:      payment.Type.String(),
				Status:    payment.Status.String(),
				Scheme:    payment.Scheme.String(),
				Asset:     payment.Asset.String(),
				Amount:    payment.Amount,
			},
		},
	}, nil
}
