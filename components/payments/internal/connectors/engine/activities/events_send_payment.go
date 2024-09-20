package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

type EventsSendPaymentRequest struct {
	Payment    models.Payment
	Adjustment models.PaymentAdjustment
}

func (a Activities) EventsSendPayment(ctx context.Context, req EventsSendPaymentRequest) error {
	return a.events.Publish(ctx, a.events.NewEventSavedPayments(req.Payment, req.Adjustment))
}

var EventsSendPaymentActivity = Activities{}.EventsSendPayment

func EventsSendPayment(ctx workflow.Context, payment models.Payment, adjustment models.PaymentAdjustment) error {
	return executeActivity(ctx, EventsSendPaymentActivity, nil, EventsSendPaymentRequest{
		Payment:    payment,
		Adjustment: adjustment,
	})
}
