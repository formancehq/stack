// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/stack/ledger/client/models/components"
)

type V2ExportLogsRequest struct {
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

func (o *V2ExportLogsRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type V2ExportLogsResponse struct {
	HTTPMeta components.HTTPMetadata `json:"-"`
}

func (o *V2ExportLogsResponse) GetHTTPMeta() components.HTTPMetadata {
	if o == nil {
		return components.HTTPMetadata{}
	}
	return o.HTTPMeta
}
