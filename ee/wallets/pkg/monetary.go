package wallet

import (
	"math/big"
)

type Monetary struct {
	Amount *big.Int `json:"amount"`
	Asset  string   `json:"asset"`
}

func NewMonetary(amount *big.Int, asset string) Monetary {
	return Monetary{
		Amount: amount,
		Asset:  asset,
	}
}
