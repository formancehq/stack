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

func UnwrapOpenAPIError(err error) *sdk.ErrorResponse {
	for err != nil {
		if err, ok := err.(*sdk.GenericOpenAPIError); ok {
			body := err.Body()
			// Actually, each api redefine errors response
			// So OpenAPI generator generate an error structure for every service
			// Manually unmarshal errorResponse allow us to handle only one ErrorResponse
			// It will be refined once the monorepo fully ready
			errResponse := api.ErrorResponse{}
			if err := json.Unmarshal(body, &errResponse); err != nil {
				return nil
			}
			if errResponse.ErrorCode != "" {
				errorCode := sdk.ErrorsEnum(errResponse.ErrorCode)
				return &sdk.ErrorResponse{
					ErrorCode:    &errorCode,
					ErrorMessage: &errResponse.ErrorMessage,
					Details:      &errResponse.Details,
				}
			}
		}

		err = errors.Unwrap(err)
	}
	return nil
}

func openApiErrorToApplicationError(err error) error {
	if err == nil {
		return nil
	}
	if err := UnwrapOpenAPIError(err); err != nil {
		return temporal.NewApplicationError(*err.ErrorMessage, string(*err.ErrorCode), err.Details)
	}
	return err
}
