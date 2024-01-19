package activities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type VoidHoldRequest struct {
	ID string `json:"id"`
}

func (a Activities) VoidHold(ctx context.Context, request VoidHoldRequest) error {
	response, err := a.client.Wallets.VoidHold(
		ctx,
		operations.VoidHoldRequest{
			HoldID: request.ID,
		},
	)
	if err != nil {
		return err
	}

	if response.WalletsErrorResponse != nil {
		return temporal.NewApplicationError(
			response.WalletsErrorResponse.ErrorMessage,
			string(response.WalletsErrorResponse.ErrorCode),
		)
	}

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var VoidHoldActivity = Activities{}.VoidHold

func VoidHold(ctx workflow.Context, id string) error {
	return executeActivity(ctx, VoidHoldActivity, nil, VoidHoldRequest{
		ID: id,
	})
}
