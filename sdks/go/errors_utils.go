package formance

import (
	"encoding/json"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pkg/errors"
)

func UnwrapOpenAPIError(err error) *ErrorResponse {
	for err != nil {
		if err, ok := err.(*GenericOpenAPIError); ok {
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
				errorCode := ErrorsEnum(errResponse.ErrorCode)
				return &ErrorResponse{
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

func ExtractOpenAPIErrorMessage(err error) error {
	if err == nil {
		return nil
	}
	if err := UnwrapOpenAPIError(err); err != nil {
		return errors.New(err.GetErrorMessage())
	}
	return err
}
