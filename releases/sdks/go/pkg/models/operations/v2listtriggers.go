// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type V2ListTriggersResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// General error
	V2Error *sdkerrors.V2Error
	// List of triggers
	V2ListTriggersResponse *shared.V2ListTriggersResponse
}

func (o *V2ListTriggersResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *V2ListTriggersResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *V2ListTriggersResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *V2ListTriggersResponse) GetV2Error() *sdkerrors.V2Error {
	if o == nil {
		return nil
	}
	return o.V2Error
}

func (o *V2ListTriggersResponse) GetV2ListTriggersResponse() *shared.V2ListTriggersResponse {
	if o == nil {
		return nil
	}
	return o.V2ListTriggersResponse
}
