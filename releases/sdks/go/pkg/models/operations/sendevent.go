// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"net/http"
)

type SendEventRequestBody struct {
	Name string `json:"name"`
}

func (o *SendEventRequestBody) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

type SendEventRequest struct {
	RequestBody *SendEventRequestBody `request:"mediaType=application/json"`
	// The instance id
	InstanceID string `pathParam:"style=simple,explode=false,name=instanceID"`
}

func (o *SendEventRequest) GetRequestBody() *SendEventRequestBody {
	if o == nil {
		return nil
	}
	return o.RequestBody
}

func (o *SendEventRequest) GetInstanceID() string {
	if o == nil {
		return ""
	}
	return o.InstanceID
}

type SendEventResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *SendEventResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *SendEventResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *SendEventResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
