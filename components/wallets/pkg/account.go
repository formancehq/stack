package wallet

import (
	"encoding/json"
	"math/big"
)

type Account struct {
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

// notes(gfyrag): hacky way to keep compatibility with ledger v1
func (a *Account) UnmarshalJSON(data []byte) error {
	type account Account
	type aux struct {
		account
		Metadata map[string]any `json:"metadata"`
	}
	v := aux{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*a = Account(v.account)
	a.Metadata = map[string]string{}
	for k, v := range v.Metadata {
		switch v := v.(type) {
		case string:
			a.Metadata[k] = v
		default:
			data, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			a.Metadata[k] = string(data)
		}
	}
	return nil
}

func (a Account) GetMetadata() map[string]string {
	return a.Metadata
}

func (a Account) GetAddress() string {
	return a.Address
}

type AccountWithVolumesAndBalances struct {
	Account
	Volumes  map[string]map[string]*big.Int `json:"volumes"`
	Balances map[string]*big.Int            `json:"balances"`
}

func (a AccountWithVolumesAndBalances) GetVolumes() map[string]map[string]*big.Int {
	return a.Volumes
}

func (a AccountWithVolumesAndBalances) GetBalances() map[string]*big.Int {
	return a.Balances
}
