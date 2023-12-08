// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type InsertConfigResponse struct {
	// Config created successfully.
	ConfigResponse *shared.ConfigResponse
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Error
	WebhooksErrorResponse *shared.WebhooksErrorResponse
}

func (o *InsertConfigResponse) GetConfigResponse() *shared.ConfigResponse {
	if o == nil {
		return nil
	}
	return o.ConfigResponse
}

func (o *InsertConfigResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *InsertConfigResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *InsertConfigResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *InsertConfigResponse) GetWebhooksErrorResponse() *shared.WebhooksErrorResponse {
	if o == nil {
		return nil
	}
	return o.WebhooksErrorResponse
}
