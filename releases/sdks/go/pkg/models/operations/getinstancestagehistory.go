// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type GetInstanceStageHistoryRequest struct {
	// The instance id
	InstanceID string `pathParam:"style=simple,explode=false,name=instanceID"`
	// The stage number
	Number int64 `pathParam:"style=simple,explode=false,name=number"`
}

func (o *GetInstanceStageHistoryRequest) GetInstanceID() string {
	if o == nil {
		return ""
	}
	return o.InstanceID
}

func (o *GetInstanceStageHistoryRequest) GetNumber() int64 {
	if o == nil {
		return 0
	}
	return o.Number
}

type GetInstanceStageHistoryResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// General error
	Error *shared.Error
	// The workflow instance stage history
	GetWorkflowInstanceHistoryStageResponse *shared.GetWorkflowInstanceHistoryStageResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *GetInstanceStageHistoryResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetInstanceStageHistoryResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}

func (o *GetInstanceStageHistoryResponse) GetGetWorkflowInstanceHistoryStageResponse() *shared.GetWorkflowInstanceHistoryStageResponse {
	if o == nil {
		return nil
	}
	return o.GetWorkflowInstanceHistoryStageResponse
}

func (o *GetInstanceStageHistoryResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetInstanceStageHistoryResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}