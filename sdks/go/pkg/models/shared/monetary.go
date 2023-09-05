// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"math/big"
)

type Monetary struct {
	// The amount of the monetary value.
	Amount *big.Int `json:"amount"`
	// The asset of the monetary value.
	Asset string `json:"asset"`
}

func (o *Monetary) GetAmount() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.Amount
}

func (o *Monetary) GetAsset() string {
	if o == nil {
		return ""
	}
	return o.Asset
}
