package models

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/gibson042/canonicaljson-go"
	"github.com/uptrace/bun"
)

type PaymentReference struct {
	Reference string
	Type      PaymentType
}

type PaymentID struct {
	PaymentReference
	ConnectorID ConnectorID
}

func (pid PaymentID) String() string {
	data, err := canonicaljson.Marshal(pid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func PaymentIDFromString(value string) (*PaymentID, error) {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return nil, err
	}
	ret := PaymentID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func MustPaymentIDFromString(value string) *PaymentID {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		panic(err)
	}
	ret := PaymentID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	return &ret
}

func (pid PaymentID) Value() (driver.Value, error) {
	return pid.String(), nil
}

func (pid *PaymentID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := PaymentIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse paymentid %s: %v", v, err)
			}

			*pid = *id
			return nil
		}
	}

	return fmt.Errorf("failed to scan paymentid: %v", value)
}

type Payment struct {
	bun.BaseModel `bun:"payments.payment"`

	ID          PaymentID   `bun:",pk,nullzero"`
	ConnectorID ConnectorID `bun:",nullzero"`
	CreatedAt   time.Time   `bun:",nullzero"`
	Reference   string
	Amount      *big.Int `bun:"type:numeric"`
	Type        PaymentType
	Status      PaymentStatus
	Scheme      PaymentScheme
	Asset       Asset

	RawData json.RawMessage

	SourceAccountID      *AccountID `bun:",nullzero"`
	DestinationAccountID *AccountID `bun:",nullzero"`

	Adjustments []*Adjustment `bun:"rel:has-many,join:id=payment_id"`
	Metadata    []*Metadata   `bun:"rel:has-many,join:id=payment_id"`
	Connector   *Connector    `bun:"rel:has-one,join:connector_id=id"`
}

type (
	PaymentType   string
	PaymentStatus string
	PaymentScheme string
	Asset         string
)

const (
	PaymentTypePayIn    PaymentType = "PAY-IN"
	PaymentTypePayOut   PaymentType = "PAYOUT"
	PaymentTypeTransfer PaymentType = "TRANSFER"
	PaymentTypeOther    PaymentType = "OTHER"
)

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSucceeded PaymentStatus = "SUCCEEDED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusOther     PaymentStatus = "OTHER"
)

const (
	PaymentSchemeUnknown PaymentScheme = "unknown"
	PaymentSchemeOther   PaymentScheme = "other"

	PaymentSchemeCardVisa       PaymentScheme = "visa"
	PaymentSchemeCardMasterCard PaymentScheme = "mastercard"
	PaymentSchemeCardAmex       PaymentScheme = "amex"
	PaymentSchemeCardDiners     PaymentScheme = "diners"
	PaymentSchemeCardDiscover   PaymentScheme = "discover"
	PaymentSchemeCardJCB        PaymentScheme = "jcb"
	PaymentSchemeCardUnionPay   PaymentScheme = "unionpay"

	PaymentSchemeSepaDebit  PaymentScheme = "sepa debit"
	PaymentSchemeSepaCredit PaymentScheme = "sepa credit"
	PaymentSchemeSepa       PaymentScheme = "sepa"

	PaymentSchemeApplePay  PaymentScheme = "apple pay"
	PaymentSchemeGooglePay PaymentScheme = "google pay"

	PaymentSchemeA2A      PaymentScheme = "a2a"
	PaymentSchemeACHDebit PaymentScheme = "ach debit"
	PaymentSchemeACH      PaymentScheme = "ach"
	PaymentSchemeRTP      PaymentScheme = "rtp"
)

func (t PaymentType) String() string {
	return string(t)
}

func (t PaymentStatus) String() string {
	return string(t)
}

func (t PaymentScheme) String() string {
	return string(t)
}

func (t Asset) String() string {
	return string(t)
}

func GetCurrencyAndPrecisionFromAsset(asset Asset) (string, int64, error) {
	parts := strings.Split(asset.String(), "/")
	if len(parts) != 2 {
		return "", 0, errors.New("invalid asset")
	}

	precision, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, errors.New("invalid asset precision")
	}

	return parts[0], precision, nil
}
