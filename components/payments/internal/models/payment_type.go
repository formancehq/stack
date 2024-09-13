package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type PaymentType int

const (
	PAYMENT_TYPE_UNKNOWN PaymentType = iota
	PAYMENT_TYPE_PAYIN
	PAYMENT_TYPE_PAYOUT
	PAYMENT_TYPE_TRANSFER
	PAYMENT_TYPE_OTHER = 100 // match grpc tag
)

func (t PaymentType) String() string {
	switch t {
	case PAYMENT_TYPE_PAYIN:
		return "PAYIN"
	case PAYMENT_TYPE_PAYOUT:
		return "PAYOUT"
	case PAYMENT_TYPE_TRANSFER:
		return "TRANSFER"
	case PAYMENT_TYPE_OTHER:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

func PaymentTypeFromString(value string) (PaymentType, error) {
	switch value {
	case "PAYIN":
		return PAYMENT_TYPE_PAYIN, nil
	case "PAYOUT":
		return PAYMENT_TYPE_PAYOUT, nil
	case "TRANSFER":
		return PAYMENT_TYPE_TRANSFER, nil
	case "OTHER":
		return PAYMENT_TYPE_OTHER, nil
	case "UNKNOWN":
		return PAYMENT_TYPE_UNKNOWN, nil
	default:
		return PAYMENT_TYPE_UNKNOWN, fmt.Errorf("unknown payment type")
	}
}

func (t PaymentType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

func (t *PaymentType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	value, err := PaymentTypeFromString(v)
	if err != nil {
		return err
	}

	*t = value

	return nil
}

func (t PaymentType) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *PaymentType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment type is nil")
	}

	s, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("failed to convert payment type")
	}

	v, ok := s.(string)
	if !ok {
		return fmt.Errorf("failed to cast payment type")
	}

	res, err := PaymentTypeFromString(v)
	if err != nil {
		return err
	}

	*t = res

	return nil
}
