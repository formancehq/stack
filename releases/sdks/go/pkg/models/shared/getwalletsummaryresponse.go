// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"math/big"
)

type GetWalletSummaryResponseData struct {
	AvailableFunds map[string]*big.Int `json:"availableFunds"`
	Balances       []BalanceWithAssets `json:"balances"`
	ExpirableFunds map[string]*big.Int `json:"expirableFunds"`
	ExpiredFunds   map[string]*big.Int `json:"expiredFunds"`
	HoldFunds      map[string]*big.Int `json:"holdFunds"`
}

func (g GetWalletSummaryResponseData) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetWalletSummaryResponseData) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetWalletSummaryResponseData) GetAvailableFunds() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.AvailableFunds
}

func (o *GetWalletSummaryResponseData) GetBalances() []BalanceWithAssets {
	if o == nil {
		return []BalanceWithAssets{}
	}
	return o.Balances
}

func (o *GetWalletSummaryResponseData) GetExpirableFunds() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.ExpirableFunds
}

func (o *GetWalletSummaryResponseData) GetExpiredFunds() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.ExpiredFunds
}

func (o *GetWalletSummaryResponseData) GetHoldFunds() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.HoldFunds
}

type GetWalletSummaryResponse struct {
	Data GetWalletSummaryResponseData `json:"data"`
}

func (o *GetWalletSummaryResponse) GetData() GetWalletSummaryResponseData {
	if o == nil {
		return GetWalletSummaryResponseData{}
	}
	return o.Data
}
