// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"math/big"
	"net/http"
)

type V2DeleteTransactionMetadataRequest struct {
	// Transaction ID.
	ID *big.Int `pathParam:"style=simple,explode=false,name=id"`
	// The key to remove.
	Key string `pathParam:"style=simple,explode=false,name=key"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
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

func (o *V2DeleteTransactionMetadataRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type V2DeleteTransactionMetadataResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Error
	V2ErrorResponse *sdkerrors.V2ErrorResponse
}

func (o *V2DeleteTransactionMetadataResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *V2DeleteTransactionMetadataResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *V2DeleteTransactionMetadataResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *V2DeleteTransactionMetadataResponse) GetV2ErrorResponse() *sdkerrors.V2ErrorResponse {
	if o == nil {
		return nil
	}
	return o.V2ErrorResponse
}
