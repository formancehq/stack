package models

import (
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/gibson042/canonicaljson-go"
)

type PaymentAdjustmentID struct {
	PaymentID
	CreatedAt time.Time
	Status    PaymentStatus
}

func (pid PaymentAdjustmentID) String() string {
	data, err := canonicaljson.Marshal(pid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func PaymentAdjustmentIDFromString(value string) (*PaymentAdjustmentID, error) {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return nil, err
	}
	ret := PaymentAdjustmentID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func MustPaymentAdjustmentIDFromString(value string) *PaymentAdjustmentID {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		panic(err)
	}
	ret := PaymentAdjustmentID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	return &ret
}

func (pid PaymentAdjustmentID) Value() (driver.Value, error) {
	return pid.String(), nil
}

func (pid *PaymentAdjustmentID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment adjustment id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := PaymentAdjustmentIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse payment adjustment id %s: %v", v, err)
			}

			*pid = *id
			return nil
		}
	}

	return fmt.Errorf("failed to scan payment adjustement id: %v", value)
}
