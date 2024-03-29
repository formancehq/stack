// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type DeletePolicyRequest struct {
	// The policy ID.
	PolicyID string `pathParam:"style=simple,explode=false,name=policyID"`
}

func (o *DeletePolicyRequest) GetPolicyID() string {
	if o == nil {
		return ""
	}
	return o.PolicyID
}

type DeletePolicyResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Error response
	ReconciliationErrorResponse *shared.ReconciliationErrorResponse
}

func (o *DeletePolicyResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeletePolicyResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeletePolicyResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *DeletePolicyResponse) GetReconciliationErrorResponse() *shared.ReconciliationErrorResponse {
	if o == nil {
		return nil
	}
	return o.ReconciliationErrorResponse
}
