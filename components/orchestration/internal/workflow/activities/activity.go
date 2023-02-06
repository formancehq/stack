package activities

import (
	"encoding/json"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/go-libs/api"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	client *sdk.APIClient
}

func New(client *sdk.APIClient) Activities {
	return Activities{
		client: client,
	}
}

func executeActivity(ctx workflow.Context, activity any, ret any, request any) error {
	if err := workflow.ExecuteActivity(ctx, activity, request).Get(ctx, ret); err != nil {
		var timeoutError *temporal.TimeoutError
		if errors.As(err, &timeoutError) {
			return errors.New(timeoutError.Message())
		}
		return err
	}
	return nil
}

func openApiErrorToApplicationError(err error) error {
	if err == nil {
		return nil
	}
	genericOpenAPIError := &sdk.GenericOpenAPIError{}
	if errors.As(err, &genericOpenAPIError) {
		body := genericOpenAPIError.Body()
		// Actually, each api redefine errors response
		// So OpenAPI generator generate an error structure for every service
		// Manually unmarshal errorResponse allow us to handle only one ErrorResponse
		// It will be refined once the monorepo fully ready
		errResponse := api.ErrorResponse{}
		if err := json.Unmarshal(body, &errResponse); err != nil {
			return nil
		}
		if errResponse.ErrorCode != "" {
			return temporal.NewApplicationError(errResponse.ErrorMessage, errResponse.ErrorCode, errResponse.Details)
		}
	}
	return err
}
