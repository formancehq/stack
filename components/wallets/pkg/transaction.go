package wallet

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type Volume struct {
	Input   *big.Int `json:"input"`
	Output  *big.Int `json:"output"`
	Balance *big.Int `json:"balance,omitempty"`
}

type Posting struct {
	Amount      *big.Int `json:"amount"`
	Asset       string   `json:"asset"`
	Destination string   `json:"destination"`
	Source      string   `json:"source"`
}

type ExpandedTransaction struct {
	Timestamp         time.Time                    `json:"timestamp"`
	Postings          []Posting                    `json:"postings"`
	Reference         *string                      `json:"reference,omitempty"`
	Metadata          metadata.Metadata            `json:"metadata"`
	ID                int64                        `json:"id"`
	PreCommitVolumes  map[string]map[string]Volume `json:"preCommitVolumes"`
	PostCommitVolumes map[string]map[string]Volume `json:"postCommitVolumes"`
}

type Transaction struct {
	ExpandedTransaction
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
