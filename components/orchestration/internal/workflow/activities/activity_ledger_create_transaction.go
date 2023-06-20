package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/mod/semver"
)

type CreateTransactionRequest struct {
	Ledger string                 `json:"ledger"`
	Data   shared.PostTransaction `json:"data"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*shared.CreateTransactionResponse, error) {
	versionsResponse, err := a.client.GetVersions(ctx)
	if err != nil {
		return nil, err
	}
	if versionsResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when getting versions", versionsResponse.StatusCode)
	}

	version := collectionutils.Filter(versionsResponse.GetVersionsResponse.Versions, func(version shared.Version) bool {
		return version.Name == "ledger"
	})[0].Version

	response, err := a.client.Ledger.
		CreateTransaction(
			ctx,
			operations.CreateTransactionRequest{
				PostTransaction: request.Data,
				Ledger:          request.Ledger,
			},
		)
	if semver.IsValid(version) && semver.Compare(version, "v2.0.0") < 0 {
		baseResponse := api.BaseResponse[[]shared.Transaction]{}
		if err := json.NewDecoder(response.RawResponse.Body).Decode(&baseResponse); err != nil {
			return nil, err
		}
		return &shared.CreateTransactionResponse{
			Data: (*baseResponse.Data)[0],
		}, nil
	} else {
		if err != nil {
			return nil, err
		}

		switch response.StatusCode {
		case http.StatusOK:
			return response.CreateTransactionResponse, nil
		default:
			if response.ErrorResponse != nil {
				return nil, temporal.NewApplicationError(
					response.ErrorResponse.ErrorMessage,
					string(response.ErrorResponse.ErrorCode),
					response.ErrorResponse.Details)
			}

			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
	}
}

var CreateTransactionActivity = Activities{}.CreateTransaction

func CreateTransaction(ctx workflow.Context, ledger string, request shared.PostTransaction) (*shared.Transaction, error) {
	tx := &shared.CreateTransactionResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data, nil
}
