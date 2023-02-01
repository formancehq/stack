package wallet

import (
	sdk "github.com/formancehq/formance-sdk-go"
)

type Transaction struct {
	sdk.Transaction
	Ledger string `json:"ledger"`
}
