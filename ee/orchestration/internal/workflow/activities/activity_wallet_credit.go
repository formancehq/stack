package activities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type CreditWalletRequest struct {
	ID   string                      `json:"id"`
	Data *shared.CreditWalletRequest `json:"data"`
}

func (a Activities) CreditWallet(ctx context.Context, request CreditWalletRequest) error {
	response, err := a.client.Wallets.CreditWallet(
		ctx,
		operations.CreditWalletRequest{
			CreditWalletRequest: request.Data,
			ID:                  request.ID,
		},
	)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return errors.New("wallet not found")
	default:
		if response.WalletsErrorResponse != nil {
			return temporal.NewApplicationError(
				response.WalletsErrorResponse.ErrorMessage,
				string(response.WalletsErrorResponse.ErrorCode),
			)
		}

		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}

var CreditWalletActivity = Activities{}.CreditWallet

func CreditWallet(ctx workflow.Context, id string, request *shared.CreditWalletRequest) error {
	return executeActivity(ctx, CreditWalletActivity, nil, CreditWalletRequest{
		ID:   id,
		Data: request,
	})
}
