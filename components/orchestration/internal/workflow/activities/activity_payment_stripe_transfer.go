package activities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StripeTransfer(ctx context.Context, request shared.ActivityStripeTransfer) error {
	activityInfo := activity.GetInfo(ctx)
	ti := shared.TransferInitiationRequest{
		Amount:               request.Amount,
		Asset:                *request.Asset,
		DestinationAccountID: *request.Destination,
		Provider:             shared.ConnectorStripe,
		Type:                 shared.TransferInitiationRequestTypeTransfer,
		Reference:            activityInfo.WorkflowExecution.ID + activityInfo.ActivityID,
		Validated:            true, // No need to validate
	}

	response, err := a.client.Payments.
		CreateTransferInitiation(
			ctx,
			ti,
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

func StripeTransfer(ctx workflow.Context, request shared.ActivityStripeTransfer) error {
	return executeActivity(ctx, StripeTransferActivity, nil, request)
}
