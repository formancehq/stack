// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type CreateClientResponse struct {
	ContentType string
	// Client created
	CreateClientResponse *shared.CreateClientResponse
	StatusCode           int
	RawResponse          *http.Response
}

func (o *CreateClientResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateClientResponse) GetCreateClientResponse() *shared.CreateClientResponse {
	if o == nil {
		return nil
	}
	return o.CreateClientResponse
}

func (o *CreateClientResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateClientResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
