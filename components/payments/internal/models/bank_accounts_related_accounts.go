package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BankAccountRelatedAccount struct {
	BankAccountID uuid.UUID   `json:"bankAccountID"`
	AccountID     AccountID   `json:"accountID"`
	ConnectorID   ConnectorID `json:"connectorID"`
	CreatedAt     time.Time   `json:"createdAt"`
}

func (b BankAccountRelatedAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		BankAccountID uuid.UUID `json:"bankAccountID"`
		AccountID     string    `json:"accountID"`
		ConnectorID   string    `json:"connectorID"`
		CreatedAt     time.Time `json:"createdAt"`
	}{
		BankAccountID: b.BankAccountID,
		AccountID:     b.AccountID.String(),
		ConnectorID:   b.ConnectorID.String(),
		CreatedAt:     b.CreatedAt,
	})
}

func (b *BankAccountRelatedAccount) UnmarshalJSON(data []byte) error {
	var aux struct {
		BankAccountID uuid.UUID `json:"bankAccountID"`
		AccountID     string    `json:"accountID"`
		ConnectorID   string    `json:"connectorID"`
		CreatedAt     time.Time `json:"createdAt"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	accountID, err := AccountIDFromString(aux.AccountID)
	if err != nil {
		return err
	}

	connectorID, err := ConnectorIDFromString(aux.ConnectorID)
	if err != nil {
		return err
	}

	b.BankAccountID = aux.BankAccountID
	b.AccountID = accountID
	b.ConnectorID = connectorID
	b.CreatedAt = aux.CreatedAt

	return nil
}
