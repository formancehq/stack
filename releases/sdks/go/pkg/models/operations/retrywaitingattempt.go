// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"net/http"
)

type RetryWaitingAttemptRequest struct {
	// Attempt ID
	AttemptID string `pathParam:"style=simple,explode=false,name=attemptId"`
}

func (o *RetryWaitingAttemptRequest) GetAttemptID() string {
	if o == nil {
		return ""
	}
	return o.AttemptID
}

type RetryWaitingAttemptResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *RetryWaitingAttemptResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *RetryWaitingAttemptResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *RetryWaitingAttemptResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
