// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

import (
	"github.com/formancehq/stack/ledger/client/internal/utils"
	"math/big"
)

type Posting struct {
	Amount      *big.Int `json:"amount"`
	Asset       string   `json:"asset"`
	Destination string   `json:"destination"`
	Source      string   `json:"source"`
}

func (p Posting) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(p, "", false)
}

func (p *Posting) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &p, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *Posting) GetAmount() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.Amount
}

func (o *Posting) GetAsset() string {
	if o == nil {
		return ""
	}
	return o.Asset
}

func (o *Posting) GetDestination() string {
	if o == nil {
		return ""
	}
	return o.Destination
}

func (o *Posting) GetSource() string {
	if o == nil {
		return ""
	}
	return o.Source
}
