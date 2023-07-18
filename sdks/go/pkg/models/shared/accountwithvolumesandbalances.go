// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"math/big"
)

type AccountWithVolumesAndBalances struct {
	Address  string                         `json:"address"`
	Balances map[string]*big.Int            `json:"balances"`
	Metadata map[string]string              `json:"metadata"`
	Type     *string                        `json:"type,omitempty"`
	Volumes  map[string]map[string]*big.Int `json:"volumes"`
}
