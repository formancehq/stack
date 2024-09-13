package models

import (
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/gibson042/canonicaljson-go"
)

type PaymentInitiationID struct {
	Reference   string
	ConnectorID ConnectorID
}

func (pid PaymentInitiationID) String() string {
	data, err := canonicaljson.Marshal(pid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func PaymentInitiationIDFromString(value string) (PaymentInitiationID, error) {
	ret := PaymentInitiationID{}
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

func MustPaymentInitiationIDFromString(value string) *PaymentInitiationID {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		panic(err)
	}
	ret := PaymentInitiationID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	return &ret
}

func (pid PaymentInitiationID) Value() (driver.Value, error) {
	return pid.String(), nil
}

func (pid *PaymentInitiationID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment initiation id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := PaymentInitiationIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse payment initiation id %s: %v", v, err)
			}

			*pid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan payment initiation id: %v", value)
}
