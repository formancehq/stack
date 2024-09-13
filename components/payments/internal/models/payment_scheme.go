package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// TODO(polo): use stringer generator
type PaymentScheme int

const (
	PAYMENT_SCHEME_UNKNOWN PaymentScheme = iota

	PAYMENT_SCHEME_CARD_VISA
	PAYMENT_SCHEME_CARD_MASTERCARD
	PAYMENT_SCHEME_CARD_AMEX
	PAYMENT_SCHEME_CARD_DINERS
	PAYMENT_SCHEME_CARD_DISCOVER
	PAYMENT_SCHEME_CARD_JCB
	PAYMENT_SCHEME_CARD_UNION_PAY
	PAYMENT_SCHEME_CARD_ALIPAY
	PAYMENT_SCHEME_CARD_CUP

	PAYMENT_SCHEME_SEPA_DEBIT
	PAYMENT_SCHEME_SEPA_CREDIT
	PAYMENT_SCHEME_SEPA

	PAYMENT_SCHEME_GOOGLE_PAY
	PAYMENT_SCHEME_APPLE_PAY

	PAYMENT_SCHEME_DOKU
	PAYMENT_SCHEME_DRAGON_PAY
	PAYMENT_SCHEME_MAESTRO
	PAYMENT_SCHEME_MOL_PAY

	PAYMENT_SCHEME_A2A
	PAYMENT_SCHEME_ACH_DEBIT
	PAYMENT_SCHEME_ACH
	PAYMENT_SCHEME_RTP

	PAYMENT_SCHEME_OTHER = 100 // match grpc tag
)

func (s PaymentScheme) String() string {
	switch s {
	case PAYMENT_SCHEME_UNKNOWN:
		return "UNKNOWN"
	case PAYMENT_SCHEME_CARD_VISA:
		return "CARD_VISA"
	case PAYMENT_SCHEME_CARD_MASTERCARD:
		return "CARD_MASTERCARD"
	case PAYMENT_SCHEME_CARD_AMEX:
		return "CARD_AMEX"
	case PAYMENT_SCHEME_CARD_DINERS:
		return "CARD_DINERS"
	case PAYMENT_SCHEME_CARD_DISCOVER:
		return "CARD_DISCOVER"
	case PAYMENT_SCHEME_CARD_JCB:
		return "CARD_JCB"
	case PAYMENT_SCHEME_CARD_UNION_PAY:
		return "CARD_UNION_PAY"
	case PAYMENT_SCHEME_CARD_ALIPAY:
		return "CARD_ALIPAY"
	case PAYMENT_SCHEME_CARD_CUP:
		return "CARD_CUP"
	case PAYMENT_SCHEME_SEPA_DEBIT:
		return "SEPA_DEBIT"
	case PAYMENT_SCHEME_SEPA_CREDIT:
		return "SEPA_CREDIT"
	case PAYMENT_SCHEME_SEPA:
		return "SEPA"
	case PAYMENT_SCHEME_GOOGLE_PAY:
		return "GOOGLE_PAY"
	case PAYMENT_SCHEME_APPLE_PAY:
		return "APPLE_PAY"
	case PAYMENT_SCHEME_DOKU:
		return "DOKU"
	case PAYMENT_SCHEME_DRAGON_PAY:
		return "DRAGON_PAY"
	case PAYMENT_SCHEME_MAESTRO:
		return "MAESTRO"
	case PAYMENT_SCHEME_MOL_PAY:
		return "MOL_PAY"
	case PAYMENT_SCHEME_A2A:
		return "A2A"
	case PAYMENT_SCHEME_ACH_DEBIT:
		return "ACH_DEBIT"
	case PAYMENT_SCHEME_ACH:
		return "ACH"
	case PAYMENT_SCHEME_RTP:
		return "RTP"
	case PAYMENT_SCHEME_OTHER:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

func PaymentSchemeFromString(value string) (PaymentScheme, error) {
	switch value {
	case "CARD_VISA":
		return PAYMENT_SCHEME_CARD_VISA, nil
	case "CARD_MASTERCARD":
		return PAYMENT_SCHEME_CARD_MASTERCARD, nil
	case "CARD_AMEX":
		return PAYMENT_SCHEME_CARD_AMEX, nil
	case "CARD_DINERS":
		return PAYMENT_SCHEME_CARD_DINERS, nil
	case "CARD_DISCOVER":
		return PAYMENT_SCHEME_CARD_DISCOVER, nil
	case "CARD_JCB":
		return PAYMENT_SCHEME_CARD_JCB, nil
	case "CARD_UNION_PAY":
		return PAYMENT_SCHEME_CARD_UNION_PAY, nil
	case "CARD_ALIPAY":
		return PAYMENT_SCHEME_CARD_ALIPAY, nil
	case "CARD_CUP":
		return PAYMENT_SCHEME_CARD_CUP, nil
	case "SEPA_DEBIT":
		return PAYMENT_SCHEME_SEPA_DEBIT, nil
	case "SEPA_CREDIT":
		return PAYMENT_SCHEME_SEPA_CREDIT, nil
	case "SEPA":
		return PAYMENT_SCHEME_SEPA, nil
	case "GOOGLE_PAY":
		return PAYMENT_SCHEME_GOOGLE_PAY, nil
	case "APPLE_PAY":
		return PAYMENT_SCHEME_APPLE_PAY, nil
	case "DOKU":
		return PAYMENT_SCHEME_DOKU, nil
	case "DRAGON_PAY":
		return PAYMENT_SCHEME_DRAGON_PAY, nil
	case "MAESTRO":
		return PAYMENT_SCHEME_MAESTRO, nil
	case "MOL_PAY":
		return PAYMENT_SCHEME_MOL_PAY, nil
	case "A2A":
		return PAYMENT_SCHEME_A2A, nil
	case "ACH_DEBIT":
		return PAYMENT_SCHEME_ACH_DEBIT, nil
	case "ACH":
		return PAYMENT_SCHEME_ACH, nil
	case "RTP":
		return PAYMENT_SCHEME_RTP, nil
	case "OTHER":
		return PAYMENT_SCHEME_OTHER, nil
	case "UNKNOWN":
		return PAYMENT_SCHEME_UNKNOWN, nil
	default:
		return PAYMENT_SCHEME_UNKNOWN, fmt.Errorf("unknown payment scheme")
	}
}

func (s PaymentScheme) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

func (s *PaymentScheme) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	value, err := PaymentSchemeFromString(v)
	if err != nil {
		return err
	}

	*s = value

	return nil
}

func (s PaymentScheme) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *PaymentScheme) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment type is nil")
	}

	st, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("failed to convert payment type")
	}

	v, ok := st.(string)
	if !ok {
		return fmt.Errorf("failed to cast payment type")
	}

	res, err := PaymentSchemeFromString(v)
	if err != nil {
		return err
	}

	*s = res

	return nil
}
