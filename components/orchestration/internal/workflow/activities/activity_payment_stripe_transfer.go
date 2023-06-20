package activities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StripeTransfer(ctx context.Context, request shared.StripeTransferRequest) error {
	response, err := a.client.Payments.
		ConnectorsStripeTransfer(
			ctx,
			request,
		)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var StripeTransferActivity = Activities{}.StripeTransfer

func StripeTransfer(ctx workflow.Context, request shared.StripeTransferRequest) error {
	return executeActivity(ctx, StripeTransferActivity, nil, request)
}
