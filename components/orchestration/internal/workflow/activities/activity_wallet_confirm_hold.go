package activities

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

type ConfirmHoldRequest struct {
	ID string `json:"id"`
}

func (a Activities) ConfirmHold(ctx context.Context, request ConfirmHoldRequest) error {
	_, err := a.client.WalletsApi.
		ConfirmHold(ctx, request.ID).
		Execute()
	if err != nil {
		return openApiErrorToApplicationError(err)
	}
	return nil
}

var ConfirmHoldActivity = Activities{}.ConfirmHold

func ConfirmHold(ctx workflow.Context, id string) error {
	return executeActivity(ctx, ConfirmHoldActivity, nil, ConfirmHoldRequest{
		ID: id,
	})
}
