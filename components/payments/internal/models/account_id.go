package models

import (
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/gibson042/canonicaljson-go"
)

type AccountID struct {
	Reference   string
	ConnectorID ConnectorID
}

func (aid *AccountID) String() string {
	if aid == nil || aid.Reference == "" {
		return ""
	}

	data, err := canonicaljson.Marshal(aid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func AccountIDFromString(value string) (AccountID, error) {
	ret := AccountID{}

	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return ret, err
	}

	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func MustAccountIDFromString(value string) AccountID {
	id, err := AccountIDFromString(value)
	if err != nil {
		panic(err)
	}
	return id
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

			*aid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan account id: %v", value)
}
