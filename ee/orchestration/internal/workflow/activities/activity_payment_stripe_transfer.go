package activities

import (
	"context"
	"math/big"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

type StripeTransferRequest struct {
	Amount      *big.Int `json:"amount,omitempty"`
	Asset       *string  `json:"asset,omitempty"`
	ConnectorID *string  `json:"connectorID,omitempty"`
	Destination *string  `json:"destination,omitempty"`
	// A set of key/value pairs that you can attach to a transfer object.
	// It can be useful for storing additional information about the transfer in a structured format.
	//
	Metadata          map[string]string `json:"metadata"`
	WaitingValidation *bool             `default:"false" json:"waitingValidation"`
}

func (a Activities) StripeTransfer(ctx context.Context, request StripeTransferRequest) error {
	validated := true
	if request.WaitingValidation != nil && *request.WaitingValidation {
		validated = false
	}

	activityInfo := activity.GetInfo(ctx)
	provider := shared.ConnectorStripe
	ti := shared.TransferInitiationRequest{
		Amount:               request.Amount,
		Asset:                *request.Asset,
		DestinationAccountID: *request.Destination,
		Description:          "Stripe Transfer",
		ConnectorID:          request.ConnectorID,
		Provider:             &provider,
		Type:                 shared.TransferInitiationRequestTypeTransfer,
		Reference:            activityInfo.WorkflowExecution.ID + activityInfo.ActivityID,
		Validated:            validated,
	}

	_, err := a.client.Payments.CreateTransferInitiation(ctx, ti)
	if err != nil {
		return err
	}

	return nil
}

var StripeTransferActivity = Activities{}.StripeTransfer

func StripeTransfer(ctx workflow.Context, request StripeTransferRequest) error {
	return executeActivity(ctx, StripeTransferActivity, nil, request)
}
