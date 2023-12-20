package dummypay

import (
	"math/big"
	"time"
)

type Kind string

const (
	KindPayment Kind = "payment"
	KindAccount Kind = "account"
)

type object struct {
	Kind    Kind     `json:"kind"`
	Payment *payment `json:"payment,omitempty"`
	Account *account `json:"account,omitempty"`
}

// payment represents a payment structure used in the generated files.
type payment struct {
	Reference            string            `json:"reference"`
	CreatedAt            time.Time         `json:"createdAt"`
	Amount               *big.Int          `json:"amount"`
	Asset                string            `json:"asset"`
	Type                 string            `json:"type"`
	Status               string            `json:"status"`
	Scheme               string            `json:"scheme"`
	SourceAccountID      string            `json:"sourceAccountId"`
	DestinationAccountID string            `json:"destinationAccountId"`
	Metadata             map[string]string `json:"metadata"`
}

type account struct {
	Reference    string            `json:"reference"`
	CreatedAt    time.Time         `json:"createdAt"`
	DefaultAsset string            `json:"defaultAsset"`
	AccountName  string            `json:"accountName"`
	Type         string            `json:"type"`
	Metadata     map[string]string `json:"metadata"`
}
