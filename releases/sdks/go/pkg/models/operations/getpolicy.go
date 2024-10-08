// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type GetPolicyRequest struct {
	// The policy ID.
	PolicyID string `pathParam:"style=simple,explode=false,name=policyID"`
}

func (o *GetPolicyRequest) GetPolicyID() string {
	if o == nil {
		return ""
	}
	return o.PolicyID
}

type GetPolicyResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// OK
	PolicyResponse *shared.PolicyResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *GetPolicyResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetPolicyResponse) GetPolicyResponse() *shared.PolicyResponse {
	if o == nil {
		return nil
	}
	return o.PolicyResponse
}

func (o *GetPolicyResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetPolicyResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
