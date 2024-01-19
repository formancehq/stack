// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type GetBalancesAggregatedRequest struct {
	// Filter balances involving given account, either as source or destination.
	Address *string `queryParam:"style=form,explode=true,name=address"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

func (o *GetBalancesAggregatedRequest) GetAddress() *string {
	if o == nil {
		return nil
	}
	return o.Address
}

func (o *GetBalancesAggregatedRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type GetBalancesAggregatedResponse struct {
	// OK
	AggregateBalancesResponse *shared.AggregateBalancesResponse
	// HTTP response content type for this operation
	ContentType string
	// Error
	ErrorResponse *sdkerrors.ErrorResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *GetBalancesAggregatedResponse) GetAggregateBalancesResponse() *shared.AggregateBalancesResponse {
	if o == nil {
		return nil
	}
	return o.AggregateBalancesResponse
}

func (o *GetBalancesAggregatedResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetBalancesAggregatedResponse) GetErrorResponse() *sdkerrors.ErrorResponse {
	if o == nil {
		return nil
	}
	return o.ErrorResponse
}

func (o *GetBalancesAggregatedResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetBalancesAggregatedResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
