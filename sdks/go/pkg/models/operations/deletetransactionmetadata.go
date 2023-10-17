// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"net/http"
)

type DeleteTransactionMetadataRequest struct {
	// Transaction ID.
	ID int64 `pathParam:"style=simple,explode=false,name=id"`
	// The key to remove.
	Key string `pathParam:"style=simple,explode=false,name=key"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

func (o *DeleteTransactionMetadataRequest) GetID() int64 {
	if o == nil {
		return 0
	}
	return o.ID
}

func (o *DeleteTransactionMetadataRequest) GetKey() string {
	if o == nil {
		return ""
	}
	return o.Key
}

func (o *DeleteTransactionMetadataRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type DeleteTransactionMetadataResponse struct {
	ContentType string
	StatusCode  int
	RawResponse *http.Response
}

func (o *DeleteTransactionMetadataResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeleteTransactionMetadataResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeleteTransactionMetadataResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
