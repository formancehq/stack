package wallet

import (
	"encoding/json"

	sdk "github.com/formancehq/formance-sdk-go"
)

type Transaction struct {
	sdk.ExpandedTransaction
	Ledger string `json:"ledger"`
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	asJSON, err := json.Marshal(t.ExpandedTransaction)
	if err != nil {
		return nil, err
	}
	asMap := make(map[string]any)
	if err := json.Unmarshal(asJSON, &asMap); err != nil {
		return nil, err
	}
	asMap["ledger"] = t.Ledger
	return json.Marshal(asMap)
}
