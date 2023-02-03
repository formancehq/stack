package activities

import (
	"context"
	"errors"
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go"
	"go.temporal.io/sdk/workflow"
)

type GetPaymentRequest struct {
	ID string `json:"id"`
}

func (a Activities) GetPayment(ctx context.Context, request GetPaymentRequest) (*sdk.PaymentResponse, error) {
	ret, httpResponse, err := a.client.PaymentsApi.
		GetPayment(ctx, request.ID).
		Execute()
	if err != nil {
		switch httpResponse.StatusCode {
		case http.StatusNotFound:
			return nil, errors.New("payment not found")
		default:
			return nil, sdk.ExtractOpenAPIErrorMessage(err)
		}
	}
	return ret, nil
}

var GetPaymentActivity = Activities{}.GetPayment

func GetPayment(ctx workflow.Context, id string) (*sdk.Payment, error) {
	ret := &sdk.PaymentResponse{}
	if err := executeActivity(ctx, GetPaymentActivity, ret, GetPaymentRequest{
		ID: id,
	}); err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
