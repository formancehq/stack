// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"math/big"
)

type AssetHolder struct {
	Assets map[string]*big.Int `json:"assets"`
}

func (o *AssetHolder) GetAssets() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.Assets
}
