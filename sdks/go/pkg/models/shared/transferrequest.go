// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"math/big"
)

type TransferRequest struct {
	Amount      *big.Int `json:"amount"`
	Asset       string   `json:"asset"`
	Destination string   `json:"destination"`
	Source      *string  `json:"source,omitempty"`
}

func (o *TransferRequest) GetAmount() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.Amount
}

func (o *TransferRequest) GetAsset() string {
	if o == nil {
		return ""
	}
	return o.Asset
}

func (o *TransferRequest) GetDestination() string {
	if o == nil {
		return ""
	}
	return o.Destination
}

func (o *TransferRequest) GetSource() *string {
	if o == nil {
		return nil
	}
	return o.Source
}
