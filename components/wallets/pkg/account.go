package wallet

import (
	"encoding/json"
	"math/big"

	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type Account struct {
	Address  string            `json:"address"`
	Metadata metadata.Metadata `json:"metadata"`
}

func (a Account) GetMetadata() metadata.Metadata {
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

// notes(gfyrag): used to shallow UnmarshalJSON of Account
func (a *AccountWithVolumesAndBalances) UnmarshalJSON(data []byte) error {
	account := Account{}
	if err := json.Unmarshal(data, &account); err != nil {
		return err
	}
	type aux struct {
		Volumes  map[string]map[string]*big.Int `json:"volumes"`
		Balances map[string]*big.Int            `json:"balances"`
	}
	v := aux{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*a = AccountWithVolumesAndBalances{
		Account:  account,
		Volumes:  v.Volumes,
		Balances: v.Balances,
	}
	return nil
}

func (a AccountWithVolumesAndBalances) GetVolumes() map[string]map[string]*big.Int {
	return a.Volumes
}

func (a AccountWithVolumesAndBalances) GetBalances() map[string]*big.Int {
	return a.Balances
}
