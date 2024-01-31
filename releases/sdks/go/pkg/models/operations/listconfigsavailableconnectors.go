// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type ListConfigsAvailableConnectorsResponse struct {
	// OK
	ConnectorsConfigsResponse *shared.ConnectorsConfigsResponse
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *ListConfigsAvailableConnectorsResponse) GetConnectorsConfigsResponse() *shared.ConnectorsConfigsResponse {
	if o == nil {
		return nil
	}
	return o.ConnectorsConfigsResponse
}

func (o *ListConfigsAvailableConnectorsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *ListConfigsAvailableConnectorsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *ListConfigsAvailableConnectorsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}