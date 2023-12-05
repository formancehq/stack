package activities

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"go.temporal.io/sdk/workflow"
)

type GetPaymentRequest struct {
	ID string `json:"id"`
}

func (a Activities) GetPayment(ctx context.Context, request GetPaymentRequest) (*shared.PaymentResponse, error) {
	response, err := a.client.Payments.GetPayment(
		ctx,
		operations.GetPaymentRequest{
			PaymentID: request.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return response.PaymentResponse, nil
	case http.StatusNotFound:
		return nil, errors.New("payment not found")
	default:
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
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
