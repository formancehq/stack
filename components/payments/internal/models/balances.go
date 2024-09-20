package models

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"

	"github.com/gibson042/canonicaljson-go"
)

type PSPBalance struct {
	// PSP account reference of the balance.
	AccountReference string

	// Balance Creation date.
	CreatedAt time.Time

	// Balance amount.
	Amount *big.Int

	// Currency. Should be in minor currencies unit.
	// For example: USD/2
	Asset string
}

type Balance struct {
	// Balance related formance account id
	AccountID AccountID `json:"accountID"`
	// Balance created at
	CreatedAt time.Time `json:"createdAt"`
	// Balance last updated at
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`

	// Currency. Should be in minor currencies unit.
	Asset string `json:"asset"`
	// Balance amount.
	Balance *big.Int `json:"balance"`
}

func (b *Balance) IdempotencyKey() string {
	var ik = struct {
		AccountID     string
		CreatedAt     int64
		LastUpdatedAt int64
	}{
		AccountID:     b.AccountID.String(),
		CreatedAt:     b.CreatedAt.UnixNano(),
		LastUpdatedAt: b.LastUpdatedAt.UnixNano(),
	}

	data, err := canonicaljson.Marshal(ik)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func (b Balance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		AccountID     string    `json:"accountID"`
		CreatedAt     time.Time `json:"createdAt"`
		LastUpdatedAt time.Time `json:"lastUpdatedAt"`

		Asset   string   `json:"asset"`
		Balance *big.Int `json:"balance"`
	}{
		AccountID:     b.AccountID.String(),
		CreatedAt:     b.CreatedAt,
		LastUpdatedAt: b.LastUpdatedAt,
		Asset:         b.Asset,
		Balance:       b.Balance,
	})
}

func (b *Balance) UnmarshalJSON(data []byte) error {
	var aux struct {
		AccountID     string    `json:"accountID"`
		CreatedAt     time.Time `json:"createdAt"`
		LastUpdatedAt time.Time `json:"lastUpdatedAt"`
		Asset         string    `json:"asset"`
		Balance       *big.Int  `json:"balance"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	accountID, err := AccountIDFromString(aux.AccountID)
	if err != nil {
		return err
	}

	b.AccountID = accountID
	b.CreatedAt = aux.CreatedAt
	b.LastUpdatedAt = aux.LastUpdatedAt
	b.Asset = aux.Asset
	b.Balance = aux.Balance

	return nil
}

type AggregatedBalance struct {
	Asset  string   `json:"asset"`
	Amount *big.Int `json:"amount"`
}

func FromPSPBalance(from PSPBalance, connectorID ConnectorID) Balance {
	return Balance{
		AccountID: AccountID{
			Reference:   from.AccountReference,
			ConnectorID: connectorID,
		},
		CreatedAt:     from.CreatedAt,
		LastUpdatedAt: from.CreatedAt,
		Asset:         from.Asset,
		Balance:       from.Amount,
	}
}

func FromPSPBalances(from []PSPBalance, connectorID ConnectorID) []Balance {
	balances := make([]Balance, 0, len(from))
	for _, b := range from {
		balances = append(balances, FromPSPBalance(b, connectorID))
	}
	return balances
}
