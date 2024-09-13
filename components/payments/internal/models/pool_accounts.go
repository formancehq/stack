package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type PoolAccounts struct {
	PoolID    uuid.UUID `json:"poolID"`
	AccountID AccountID `json:"accountID"`
}

func (p PoolAccounts) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		PoolID    uuid.UUID `json:"poolID"`
		AccountID string    `json:"accountID"`
	}{
		PoolID:    p.PoolID,
		AccountID: p.AccountID.String(),
	})
}

func (p *PoolAccounts) UnmarshalJSON(data []byte) error {
	var aux struct {
		PoolID    uuid.UUID `json:"poolID"`
		AccountID string    `json:"accountID"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	accountID, err := AccountIDFromString(aux.AccountID)
	if err != nil {
		return err
	}

	p.PoolID = aux.PoolID
	p.AccountID = accountID

	return nil
}
