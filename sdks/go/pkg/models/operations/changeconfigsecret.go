// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type ChangeConfigSecretRequest struct {
	ConfigChangeSecret *shared.ConfigChangeSecret `request:"mediaType=application/json"`
	// Config ID
	ID string `pathParam:"style=simple,explode=false,name=id"`
}

type ChangeConfigSecretResponse struct {
	// Secret successfully changed.
	ConfigResponse *shared.ConfigResponse
	ContentType    string
	StatusCode     int
	RawResponse    *http.Response
	// Error
	WebhooksErrorResponse *shared.WebhooksErrorResponse
}
