// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type V2ListTriggersOccurrencesRequest struct {
	// The trigger id
	TriggerID string `pathParam:"style=simple,explode=false,name=triggerID"`
}

func (o *V2ListTriggersOccurrencesRequest) GetTriggerID() string {
	if o == nil {
		return ""
	}
	return o.TriggerID
}

type V2ListTriggersOccurrencesResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// General error
	V2Error *sdkerrors.V2Error
	// List of triggers occurrences
	V2ListTriggersOccurrencesResponse *shared.V2ListTriggersOccurrencesResponse
}

func (o *V2ListTriggersOccurrencesResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *V2ListTriggersOccurrencesResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *V2ListTriggersOccurrencesResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *V2ListTriggersOccurrencesResponse) GetV2Error() *sdkerrors.V2Error {
	if o == nil {
		return nil
	}
	return o.V2Error
}

func (o *V2ListTriggersOccurrencesResponse) GetV2ListTriggersOccurrencesResponse() *shared.V2ListTriggersOccurrencesResponse {
	if o == nil {
		return nil
	}
	return o.V2ListTriggersOccurrencesResponse
}
