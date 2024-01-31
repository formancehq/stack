// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"math/big"
)

type Monetary struct {
	// The amount of the monetary value.
	Amount *big.Int `json:"amount"`
	// The asset of the monetary value.
	Asset string `json:"asset"`
}

func (m Monetary) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(m, "", false)
}

func (m *Monetary) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &m, "", false, false); err != nil {
		return err
	}
	return nil
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