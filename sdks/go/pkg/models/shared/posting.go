// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"math/big"
)

type Posting struct {
	Amount      *big.Int `json:"amount"`
	Asset       string   `json:"asset"`
	Destination string   `json:"destination"`
	Source      string   `json:"source"`
}
