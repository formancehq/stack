// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type ListTriggersResponse struct {
	ContentType string
	// General error
	Error *shared.Error
	// List of triggers
	ListTriggersResponse *shared.ListTriggersResponse
	StatusCode           int
	RawResponse          *http.Response
}
