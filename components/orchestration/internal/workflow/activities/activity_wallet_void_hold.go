package activities

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

type VoidHoldRequest struct {
	ID string `json:"id"`
}

func (a Activities) VoidHold(ctx context.Context, request VoidHoldRequest) error {
	_, err := a.client.WalletsApi.
		VoidHold(ctx, request.ID).
		Execute()
	if err != nil {
		return openApiErrorToApplicationError(err)
	}
	return nil
}

var VoidHoldActivity = Activities{}.VoidHold

func VoidHold(ctx workflow.Context, id string) error {
	return executeActivity(ctx, VoidHoldActivity, nil, VoidHoldRequest{
		ID: id,
	})
}
