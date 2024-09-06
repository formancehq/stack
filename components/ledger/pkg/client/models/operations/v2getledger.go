// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/stack/ledger/client/models/components"
)

type V2GetLedgerRequest struct {
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

func (o *V2GetLedgerRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type V2GetLedgerResponse struct {
	HTTPMeta components.HTTPMetadata `json:"-"`
	// OK
	V2GetLedgerResponse *components.V2GetLedgerResponse
}

func (o *V2GetLedgerResponse) GetHTTPMeta() components.HTTPMetadata {
	if o == nil {
		return components.HTTPMetadata{}
	}
	return o.HTTPMeta
}

func (o *V2GetLedgerResponse) GetV2GetLedgerResponse() *components.V2GetLedgerResponse {
	if o == nil {
		return nil
	}
	return o.V2GetLedgerResponse
}
