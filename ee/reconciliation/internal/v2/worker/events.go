package worker

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

// TODO(polo): use events corresponding to the current version of payments
// and ledger. We should fetch the information through the sdk and use the
// proper version of the event in the future

const (
	EventTypeCommittedTransactions = "COMMITTED_TRANSACTIONS"
	EventTypeRevertedTransaction   = "REVERTED_TRANSACTION"
	EventTypeSavedPayments         = "SAVED_PAYMENT"
)

type paymentMessagePayload struct {
	ID                   string            `json:"id"`
	Reference            string            `json:"reference"`
	CreatedAt            time.Time         `json:"createdAt"`
	ConnectorID          string            `json:"connectorId"`
	Provider             string            `json:"provider"`
	Type                 string            `json:"type"`
	Status               string            `json:"status"`
	Scheme               string            `json:"scheme"`
	Asset                string            `json:"asset"`
	SourceAccountID      string            `json:"sourceAccountId,omitempty"`
	DestinationAccountID string            `json:"destinationAccountId,omitempty"`
	Links                []api.Link        `json:"links"`
	RawData              json.RawMessage   `json:"rawData"`
	InitialAmount        *big.Int          `json:"initialAmount"`
	Amount               *big.Int          `json:"amount"`
	Metadata             map[string]string `json:"metadata"`
}

type Posting struct {
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	Amount      *big.Int `json:"amount"`
	Asset       string   `json:"asset"`
}

type Postings []Posting

type TransactionData struct {
	Postings  Postings          `json:"postings"`
	Metadata  metadata.Metadata `json:"metadata"`
	Timestamp time.Time         `json:"timestamp"`
	Reference string            `json:"reference,omitempty"`
}

type Transaction struct {
	TransactionData
	ID       *big.Int `json:"id"`
	Reverted bool     `json:"reverted"`
}

type CommittedTransactions struct {
	Ledger          string                       `json:"ledger"`
	Transactions    []Transaction                `json:"transactions"`
	AccountMetadata map[string]metadata.Metadata `json:"accountMetadata"`
}

type RevertedTransaction struct {
	Ledger              string      `json:"ledger"`
	RevertedTransaction Transaction `json:"revertedTransaction"`
	RevertTransaction   Transaction `json:"revertTransaction"`
}
