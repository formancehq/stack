// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"net/http"
)

type GetBalancesAggregatedRequest struct {
	// Filter balances involving given account, either as source or destination.
	Address *string `queryParam:"style=form,explode=true,name=address"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

type GetBalancesAggregatedResponse struct {
	// OK
	AggregateBalancesResponse *shared.AggregateBalancesResponse
	ContentType               string
	// Error
	ErrorResponse *shared.ErrorResponse
	StatusCode    int
	RawResponse   *http.Response
}
