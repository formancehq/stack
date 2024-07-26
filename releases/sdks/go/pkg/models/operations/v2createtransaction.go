// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type V2CreateTransactionRequest struct {
	// Use an idempotency key
	IdempotencyKey *string `header:"style=simple,explode=false,name=Idempotency-Key"`
	// The request body must contain at least one of the following objects:
	//   - `postings`: suitable for simple transactions
	//   - `script`: enabling more complex transactions with Numscript
	//
	V2PostTransaction shared.V2PostTransaction `request:"mediaType=application/json"`
	// Set the dryRun mode. dry run mode doesn't add the logs to the database or publish a message to the message broker.
	DryRun *bool `queryParam:"style=form,explode=true,name=dryRun"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

func (o *V2CreateTransactionRequest) GetIdempotencyKey() *string {
	if o == nil {
		return nil
	}
	return o.IdempotencyKey
}

func (o *V2CreateTransactionRequest) GetV2PostTransaction() shared.V2PostTransaction {
	if o == nil {
		return shared.V2PostTransaction{}
	}
	return o.V2PostTransaction
}

func (o *V2CreateTransactionRequest) GetDryRun() *bool {
	if o == nil {
		return nil
	}
	return o.DryRun
}

func (o *V2CreateTransactionRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type V2CreateTransactionResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// OK
	V2CreateTransactionResponse *shared.V2CreateTransactionResponse
}

func (o *V2CreateTransactionResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *V2CreateTransactionResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *V2CreateTransactionResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *V2CreateTransactionResponse) GetV2CreateTransactionResponse() *shared.V2CreateTransactionResponse {
	if o == nil {
		return nil
	}
	return o.V2CreateTransactionResponse
}
