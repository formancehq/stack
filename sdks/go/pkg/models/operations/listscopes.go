// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type ListScopesResponse struct {
	ContentType string
	// List of scopes
	ListScopesResponse *shared.ListScopesResponse
	StatusCode         int
	RawResponse        *http.Response
}

func (o *ListScopesResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *ListScopesResponse) GetListScopesResponse() *shared.ListScopesResponse {
	if o == nil {
		return nil
	}
	return o.ListScopesResponse
}

func (o *ListScopesResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *ListScopesResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
