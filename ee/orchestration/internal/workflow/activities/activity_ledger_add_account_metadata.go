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

	body := make(map[string]any)
	for k, v := range request.Metadata {
		body[k] = v
	}

	_, err := a.client.Ledger.AddMetadataToAccount(
		ctx,
		operations.AddMetadataToAccountRequest{
			RequestBody: body,
			Address:     request.Account,
			Ledger:      request.Ledger,
		},
	)
	if err != nil {
		switch err := err.(type) {
		case *sdkerrors.ErrorResponse:
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
