package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type PaymentStatus int

const (
	PAYMENT_STATUS_UNKNOWN PaymentStatus = iota
	PAYMENT_STATUS_PENDING
	PAYMENT_STATUS_SUCCEEDED
	PAYMENT_STATUS_CANCELLED
	PAYMENT_STATUS_FAILED
	PAYMENT_STATUS_EXPIRED
	PAYMENT_STATUS_REFUNDED
	PAYMENT_STATUS_REFUNDED_FAILURE
	PAYMENT_STATUS_DISPUTE
	PAYMENT_STATUS_DISPUTE_WON
	PAYMENT_STATUS_DISPUTE_LOST
	PAYMENT_STATUS_OTHER = 100 // match grpc tag
)

func (t PaymentStatus) String() string {
	switch t {
	case PAYMENT_STATUS_UNKNOWN:
		return "UNKNOWN"
	case PAYMENT_STATUS_PENDING:
		return "PENDING"
	case PAYMENT_STATUS_SUCCEEDED:
		return "SUCCEEDED"
	case PAYMENT_STATUS_CANCELLED:
		return "CANCELLED"
	case PAYMENT_STATUS_FAILED:
		return "FAILED"
	case PAYMENT_STATUS_EXPIRED:
		return "EXPIRED"
	case PAYMENT_STATUS_REFUNDED:
		return "REFUNDED"
	case PAYMENT_STATUS_REFUNDED_FAILURE:
		return "REFUNDED_FAILURE"
	case PAYMENT_STATUS_DISPUTE:
		return "DISPUTE"
	case PAYMENT_STATUS_DISPUTE_WON:
		return "DISPUTE_WON"
	case PAYMENT_STATUS_DISPUTE_LOST:
		return "DISPUTE_LOST"
	case PAYMENT_STATUS_OTHER:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

func PaymentStatusFromString(value string) (PaymentStatus, error) {
	switch value {
	case "PENDING":
		return PAYMENT_STATUS_PENDING, nil
	case "SUCCEEDED":
		return PAYMENT_STATUS_SUCCEEDED, nil
	case "CANCELLED":
		return PAYMENT_STATUS_CANCELLED, nil
	case "FAILED":
		return PAYMENT_STATUS_FAILED, nil
	case "EXPIRED":
		return PAYMENT_STATUS_EXPIRED, nil
	case "REFUNDED":
		return PAYMENT_STATUS_REFUNDED, nil
	case "REFUNDED_FAILURE":
		return PAYMENT_STATUS_REFUNDED_FAILURE, nil
	case "DISPUTE":
		return PAYMENT_STATUS_DISPUTE, nil
	case "DISPUTE_WON":
		return PAYMENT_STATUS_DISPUTE_WON, nil
	case "DISPUTE_LOST":
		return PAYMENT_STATUS_DISPUTE_LOST, nil
	case "OTHER":
		return PAYMENT_STATUS_OTHER, nil
	case "UNKNOWN":
		return PAYMENT_STATUS_UNKNOWN, nil
	default:
		return PAYMENT_STATUS_UNKNOWN, fmt.Errorf("unknown payment status")
	}
}

func (t PaymentStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

func (t *PaymentStatus) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	value, err := PaymentStatusFromString(v)
	if err != nil {
		return err
	}

	*t = value

	return nil
}

func (t PaymentStatus) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment status is nil")
	}

	s, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("failed to convert payment status")
	}

	v, ok := s.(string)
	if !ok {
		return fmt.Errorf("failed to cast payment status")
	}

	res, err := PaymentStatusFromString(v)
	if err != nil {
		return err
	}

	*t = res

	return nil
}
