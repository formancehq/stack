package wallet

import (
	"math/big"
)

type Account struct {
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
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
