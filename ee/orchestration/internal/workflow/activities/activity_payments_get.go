package activities

import (
	"context"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"go.temporal.io/sdk/workflow"
)

type GetPaymentRequest struct {
	ID string `json:"id"`
}

func (a Activities) GetPayment(ctx context.Context, request GetPaymentRequest) (*shared.PaymentResponse, error) {
	response, err := a.client.Payments.V1.GetPayment(
		ctx,
		operations.GetPaymentRequest{
			PaymentID: request.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	return response.PaymentResponse, nil
}

var GetPaymentActivity = Activities{}.GetPayment

func GetPayment(ctx workflow.Context, id string) (*shared.Payment, error) {
	ret := &shared.PaymentResponse{}
	if err := executeActivity(ctx, GetPaymentActivity, ret, GetPaymentRequest{
		ID: id,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
