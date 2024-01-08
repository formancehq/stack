// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"math/big"
	"net/http"
)

type V2RevertTransactionRequest struct {
	// Force revert
	Force *bool `queryParam:"style=form,explode=true,name=force"`
	// Transaction ID.
	ID *big.Int `pathParam:"style=simple,explode=false,name=id"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

func (v V2RevertTransactionRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(v, "", false)
}

func (v *V2RevertTransactionRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &v, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *V2RevertTransactionRequest) GetForce() *bool {
	if o == nil {
		return nil
	}
	return o.Force
}

func (o *V2RevertTransactionRequest) GetID() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.ID
}

func (o *V2RevertTransactionRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type V2RevertTransactionResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Error
	V2ErrorResponse *sdkerrors.V2ErrorResponse
	// OK
	V2RevertTransactionResponse *shared.V2RevertTransactionResponse
}

func (o *V2RevertTransactionResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *V2RevertTransactionResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *V2RevertTransactionResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *V2RevertTransactionResponse) GetV2ErrorResponse() *sdkerrors.V2ErrorResponse {
	if o == nil {
		return nil
	}
	return o.V2ErrorResponse
}

func (o *V2RevertTransactionResponse) GetV2RevertTransactionResponse() *shared.V2RevertTransactionResponse {
	if o == nil {
		return nil
	}
	return o.V2RevertTransactionResponse
}
