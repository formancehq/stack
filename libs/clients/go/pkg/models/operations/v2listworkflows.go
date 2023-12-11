// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type V2ListWorkflowsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// General error
	V2Error *shared.V2Error
	// List of workflows
	V2ListWorkflowsResponse *shared.V2ListWorkflowsResponse
}

func (o *V2ListWorkflowsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *V2ListWorkflowsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *V2ListWorkflowsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *V2ListWorkflowsResponse) GetV2Error() *shared.V2Error {
	if o == nil {
		return nil
	}
	return o.V2Error
}

func (o *V2ListWorkflowsResponse) GetV2ListWorkflowsResponse() *shared.V2ListWorkflowsResponse {
	if o == nil {
		return nil
	}
	return o.V2ListWorkflowsResponse
}
