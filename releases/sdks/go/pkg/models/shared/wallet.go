// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"time"
)

type WalletBalances struct {
	Main AssetHolder `json:"main"`
}

func (o *WalletBalances) GetMain() AssetHolder {
	if o == nil {
		return AssetHolder{}
	}
	return o.Main
}

type Wallet struct {
	Balances  *WalletBalances `json:"balances,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
	// The unique ID of the wallet.
	ID     string `json:"id"`
	Ledger string `json:"ledger"`
	// Metadata associated with the wallet.
	Metadata map[string]string `json:"metadata"`
	Name     string            `json:"name"`
}

func (w Wallet) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(w, "", false)
}

func (w *Wallet) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &w, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *Wallet) GetBalances() *WalletBalances {
	if o == nil {
		return nil
	}
	return o.Balances
}

func (o *Wallet) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *Wallet) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *Wallet) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

func (o *Wallet) GetMetadata() map[string]string {
	if o == nil {
		return map[string]string{}
	}
	return o.Metadata
}

func (o *Wallet) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}
