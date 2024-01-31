// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type GetReconciliationRequest struct {
	// The reconciliation ID.
	ReconciliationID string `pathParam:"style=simple,explode=false,name=reconciliationID"`
}

func (o *GetReconciliationRequest) GetReconciliationID() string {
	if o == nil {
		return ""
	}
	return o.ReconciliationID
}

type GetReconciliationResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// OK
	ReconciliationResponse *shared.ReconciliationResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Error response
	ReconciliationErrorResponse *shared.ReconciliationErrorResponse
}

func (o *GetReconciliationResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetReconciliationResponse) GetReconciliationResponse() *shared.ReconciliationResponse {
	if o == nil {
		return nil
	}
	return o.ReconciliationResponse
}

func (o *GetReconciliationResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetReconciliationResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetReconciliationResponse) GetReconciliationErrorResponse() *shared.ReconciliationErrorResponse {
	if o == nil {
		return nil
	}
	return o.ReconciliationErrorResponse
}