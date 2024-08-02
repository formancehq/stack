package activities

import (
	"context"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type AddAccountMetadataRequest struct {
	Ledger   string            `json:"ledger"`
	Account  string            `json:"account"`
	Metadata map[string]string `json:"metadata"`
}

func (a Activities) AddAccountMetadata(ctx context.Context, request AddAccountMetadataRequest) error {
	_, err := a.client.Ledger.V2.AddMetadataToAccount(
		ctx,
		operations.V2AddMetadataToAccountRequest{
			RequestBody:    request.Metadata,
			Address:        request.Account,
			Ledger:         request.Ledger,
			IdempotencyKey: getIK(ctx),
		},
	)
	if err != nil {
		switch err := err.(type) {
		case *sdkerrors.V2ErrorResponse:
			return temporal.NewApplicationError(err.ErrorMessage, string(err.ErrorCode), err.Details)
		default:
			return err
		}
	}

	return nil
}

var AddAccountMetadataActivity = Activities{}.AddAccountMetadata

func AddAccountMetadata(ctx workflow.Context, request AddAccountMetadataRequest) error {
	return executeActivity(ctx, AddAccountMetadataActivity, nil, request)
}
