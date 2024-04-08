package activities

import (
	"context"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"go.temporal.io/sdk/workflow"
)

type VoidHoldRequest struct {
	ID string `json:"id"`
}

func (a Activities) VoidHold(ctx context.Context, request VoidHoldRequest) error {
	_, err := a.client.Wallets.VoidHold(
		ctx,
		operations.VoidHoldRequest{
			HoldID: request.ID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

var VoidHoldActivity = Activities{}.VoidHold

func VoidHold(ctx workflow.Context, id string) error {
	return executeActivity(ctx, VoidHoldActivity, nil, VoidHoldRequest{
		ID: id,
	})
}
