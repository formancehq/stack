// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/stack/ledger/client/internal/utils"
	"github.com/formancehq/stack/ledger/client/models/components"
	"math/big"
)

type V2DeleteTransactionMetadataRequest struct {
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
	// Transaction ID.
	ID *big.Int `pathParam:"style=simple,explode=false,name=id"`
	// The key to remove.
	Key string `pathParam:"style=simple,explode=false,name=key"`
}

func (v V2DeleteTransactionMetadataRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(v, "", false)
}

func (v *V2DeleteTransactionMetadataRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &v, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *V2DeleteTransactionMetadataRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

func (o *V2DeleteTransactionMetadataRequest) GetID() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.ID
}

func (o *V2DeleteTransactionMetadataRequest) GetKey() string {
	if o == nil {
		return ""
	}
	return o.Key
}

type V2DeleteTransactionMetadataResponse struct {
	HTTPMeta components.HTTPMetadata `json:"-"`
}

func (o *V2DeleteTransactionMetadataResponse) GetHTTPMeta() components.HTTPMetadata {
	if o == nil {
		return components.HTTPMetadata{}
	}
	return o.HTTPMeta
}
