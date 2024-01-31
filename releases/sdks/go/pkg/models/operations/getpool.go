// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type GetPoolRequest struct {
	// The pool ID.
	PoolID string `pathParam:"style=simple,explode=false,name=poolId"`
}

func (o *GetPoolRequest) GetPoolID() string {
	if o == nil {
		return ""
	}
	return o.PoolID
}

type GetPoolResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// OK
	PoolResponse *shared.PoolResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *GetPoolResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetPoolResponse) GetPoolResponse() *shared.PoolResponse {
	if o == nil {
		return nil
	}
	return o.PoolResponse
}

func (o *GetPoolResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetPoolResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}