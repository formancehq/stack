package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Capability int

const (
	CAPABILITY_FETCH_UNKNOWN Capability = iota
	CAPABILITY_FETCH_ACCOUNTS
	CAPABILITY_FETCH_BALANCES
	CAPABILITY_FETCH_EXTERNAL_ACCOUNTS
	CAPABILITY_FETCH_PAYMENTS
	CAPABILITY_FETCH_OTHERS
	CAPABILITY_WEBHOOKS
	CAPABILITY_CREATION_BANK_ACCOUNT
	CAPABILITY_CREATION_PAYMENT
)

func (t Capability) String() string {
	switch t {
	case CAPABILITY_FETCH_ACCOUNTS:
		return "FETCH_ACCOUNTS"
	case CAPABILITY_FETCH_EXTERNAL_ACCOUNTS:
		return "FETCH_EXTERNAL_ACCOUNTS"
	case CAPABILITY_FETCH_PAYMENTS:
		return "FETCH_PAYMENTS"
	case CAPABILITY_FETCH_OTHERS:
		return "FETCH_OTHERS"
	case CAPABILITY_WEBHOOKS:
		return "WEBHOOKS"
	case CAPABILITY_CREATION_BANK_ACCOUNT:
		return "CREATION_BANK_ACCOUNT"
	case CAPABILITY_CREATION_PAYMENT:
		return "CREATION_PAYMENT"
	default:
		return "UNKNOWN"
	}
}

func (t Capability) Value() (driver.Value, error) {
	switch t {
	case CAPABILITY_FETCH_ACCOUNTS:
		return "FETCH_ACCOUNTS", nil
	case CAPABILITY_FETCH_EXTERNAL_ACCOUNTS:
		return "FETCH_EXTERNAL_ACCOUNTS", nil
	case CAPABILITY_FETCH_PAYMENTS:
		return "FETCH_PAYMENTS", nil
	case CAPABILITY_FETCH_OTHERS:
		return "FETCH_OTHERS", nil
	case CAPABILITY_WEBHOOKS:
		return "WEBHOOKS", nil
	case CAPABILITY_CREATION_BANK_ACCOUNT:
		return "CREATION_BANK_ACCOUNT", nil
	case CAPABILITY_CREATION_PAYMENT:
		return "CREATION_PAYMENT", nil
	default:
		return nil, fmt.Errorf("unknown capability")
	}
}

func (t *Capability) Scan(value interface{}) error {
	if value == nil {
		return errors.New("capability is nil")
	}

	s, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("failed to convert capability")
	}

	v, ok := s.(string)
	if !ok {
		return fmt.Errorf("failed to cast capability")
	}

	switch v {
	case "FETCH_ACCOUNTS":
		*t = CAPABILITY_FETCH_ACCOUNTS
	case "FETCH_EXTERNAL_ACCOUNTS":
		*t = CAPABILITY_FETCH_EXTERNAL_ACCOUNTS
	case "FETCH_PAYMENTS":
		*t = CAPABILITY_FETCH_PAYMENTS
	case "FETCH_OTHERS":
		*t = CAPABILITY_FETCH_OTHERS
	case "WEBHOOKS":
		*t = CAPABILITY_WEBHOOKS
	case "CREATION_BANK_ACCOUNT":
		*t = CAPABILITY_CREATION_BANK_ACCOUNT
	case "CREATION_PAYMENT":
		*t = CAPABILITY_CREATION_PAYMENT
	default:
		return fmt.Errorf("unknown capability")
	}

	return nil
}
