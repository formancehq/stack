package activities

import (
	"context"

	sdk "github.com/formancehq/formance-sdk-go"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StripeTransfer(ctx context.Context, request sdk.StripeTransferRequest) error {
	_, _, err := a.client.PaymentsApi.
		ConnectorsStripeTransfer(ctx).
		StripeTransferRequest(request).
		Execute()
	return openApiErrorToApplicationError(err)
}

var StripeTransferActivity = Activities{}.StripeTransfer

func StripeTransfer(ctx workflow.Context, request sdk.StripeTransferRequest) error {
	return executeActivity(ctx, StripeTransferActivity, nil, request)
}
