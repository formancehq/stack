package models

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gibson042/canonicaljson-go"
	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"accounts.account"`

	ID              AccountID `bun:",pk,nullzero"`
	CreatedAt       time.Time `bun:",nullzero"`
	Reference       string
	Provider        string
	DefaultCurrency string // Is optional and default to ''
	AccountName     string // Is optional and default to ''
	Type            AccountType

	RawData json.RawMessage

	Payments []*Payment `bun:"rel:has-many,join:id=account_id"`
}

type AccountType string

const (
	AccountTypeUnknown AccountType = "UNKNOWN"
	// Refers to an account that is internal to the psp, an account that we
	// can actually fetch the balance.
	AccountTypeInternal AccountType = "INTERNAL"
	// Refers to an external accounts such as user's bank accounts.
	AccountTypeExternal AccountType = "EXTERNAL"
)

type AccountID struct {
	Reference string
	Provider  ConnectorProvider
}

func (aid AccountID) String() string {
	data, err := canonicaljson.Marshal(aid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(data)
}

func AccountIDFromString(value string) (*AccountID, error) {
	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	ret := AccountID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (aid AccountID) Value() (driver.Value, error) {
	return aid.String(), nil
}

func (aid *AccountID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("account id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := AccountIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse account id %s: %v", v, err)
			}

			*aid = *id
			return nil
		}
	}

	return fmt.Errorf("failed to scan account id: %v", value)
}
