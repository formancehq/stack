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
	provider := shared.ConnectorStripe
	ti := shared.TransferInitiationRequest{
		Amount:               request.Amount,
		Asset:                *request.Asset,
		DestinationAccountID: *request.Destination,
		ConnectorID:          request.ConnectorID,
		Provider:             &provider,
		Type:                 shared.TransferInitiationRequestTypeTransfer,
		Reference:            activityInfo.WorkflowExecution.ID + activityInfo.ActivityID,
		Validated:            true, // No need to validate
	}

	fmt.Println("TATTA 2")
	response, err := a.client.Payments.
		CreateTransferInitiation(
			ctx,
			ti,
		)
	if err != nil {
		fmt.Println("TATTA 3")
		return err
	}
	fmt.Println("TATTA 4")

	switch response.StatusCode {
	case http.StatusOK:
		fmt.Println("TATTA 5")
		return nil
	default:
		fmt.Println("TATTA 6")
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var StripeTransferActivity = Activities{}.StripeTransfer

func StripeTransfer(ctx workflow.Context, request shared.ActivityStripeTransfer) error {
	fmt.Println("TATTA")
	return executeActivity(ctx, StripeTransferActivity, nil, request)
}
