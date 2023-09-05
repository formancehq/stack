// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

type PaymentRaw struct {
}

type PaymentScheme string

const (
	PaymentSchemeVisa       PaymentScheme = "visa"
	PaymentSchemeMastercard PaymentScheme = "mastercard"
	PaymentSchemeAmex       PaymentScheme = "amex"
	PaymentSchemeDiners     PaymentScheme = "diners"
	PaymentSchemeDiscover   PaymentScheme = "discover"
	PaymentSchemeJcb        PaymentScheme = "jcb"
	PaymentSchemeUnionpay   PaymentScheme = "unionpay"
	PaymentSchemeSepaDebit  PaymentScheme = "sepa debit"
	PaymentSchemeSepaCredit PaymentScheme = "sepa credit"
	PaymentSchemeSepa       PaymentScheme = "sepa"
	PaymentSchemeApplePay   PaymentScheme = "apple pay"
	PaymentSchemeGooglePay  PaymentScheme = "google pay"
	PaymentSchemeA2a        PaymentScheme = "a2a"
	PaymentSchemeAchDebit   PaymentScheme = "ach debit"
	PaymentSchemeAch        PaymentScheme = "ach"
	PaymentSchemeRtp        PaymentScheme = "rtp"
	PaymentSchemeUnknown    PaymentScheme = "unknown"
	PaymentSchemeOther      PaymentScheme = "other"
)

func (e PaymentScheme) ToPointer() *PaymentScheme {
	return &e
}

func (e *PaymentScheme) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "visa":
		fallthrough
	case "mastercard":
		fallthrough
	case "amex":
		fallthrough
	case "diners":
		fallthrough
	case "discover":
		fallthrough
	case "jcb":
		fallthrough
	case "unionpay":
		fallthrough
	case "sepa debit":
		fallthrough
	case "sepa credit":
		fallthrough
	case "sepa":
		fallthrough
	case "apple pay":
		fallthrough
	case "google pay":
		fallthrough
	case "a2a":
		fallthrough
	case "ach debit":
		fallthrough
	case "ach":
		fallthrough
	case "rtp":
		fallthrough
	case "unknown":
		fallthrough
	case "other":
		*e = PaymentScheme(v)
		return nil
	default:
		return fmt.Errorf("invalid value for PaymentScheme: %v", v)
	}
}

type PaymentType string

const (
	PaymentTypePayIn    PaymentType = "PAY-IN"
	PaymentTypePayout   PaymentType = "PAYOUT"
	PaymentTypeTransfer PaymentType = "TRANSFER"
	PaymentTypeOther    PaymentType = "OTHER"
)

func (e PaymentType) ToPointer() *PaymentType {
	return &e
}

func (e *PaymentType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "PAY-IN":
		fallthrough
	case "PAYOUT":
		fallthrough
	case "TRANSFER":
		fallthrough
	case "OTHER":
		*e = PaymentType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for PaymentType: %v", v)
	}
}

type Payment struct {
	Adjustments          []PaymentAdjustment `json:"adjustments"`
	Asset                string              `json:"asset"`
	CreatedAt            time.Time           `json:"createdAt"`
	DestinationAccountID string              `json:"destinationAccountID"`
	ID                   string              `json:"id"`
	InitialAmount        *big.Int            `json:"initialAmount"`
	Metadata             PaymentMetadata     `json:"metadata"`
	Provider             Connector           `json:"provider"`
	Raw                  PaymentRaw          `json:"raw"`
	Reference            string              `json:"reference"`
	Scheme               PaymentScheme       `json:"scheme"`
	SourceAccountID      string              `json:"sourceAccountID"`
	Status               PaymentStatus       `json:"status"`
	Type                 PaymentType         `json:"type"`
}

func (o *Payment) GetAdjustments() []PaymentAdjustment {
	if o == nil {
		return []PaymentAdjustment{}
	}
	return o.Adjustments
}

func (o *Payment) GetAsset() string {
	if o == nil {
		return ""
	}
	return o.Asset
}

func (o *Payment) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *Payment) GetDestinationAccountID() string {
	if o == nil {
		return ""
	}
	return o.DestinationAccountID
}

func (o *Payment) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *Payment) GetInitialAmount() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.InitialAmount
}

func (o *Payment) GetMetadata() PaymentMetadata {
	if o == nil {
		return PaymentMetadata{}
	}
	return o.Metadata
}

func (o *Payment) GetProvider() Connector {
	if o == nil {
		return Connector("")
	}
	return o.Provider
}

func (o *Payment) GetRaw() PaymentRaw {
	if o == nil {
		return PaymentRaw{}
	}
	return o.Raw
}

func (o *Payment) GetReference() string {
	if o == nil {
		return ""
	}
	return o.Reference
}

func (o *Payment) GetScheme() PaymentScheme {
	if o == nil {
		return PaymentScheme("")
	}
	return o.Scheme
}

func (o *Payment) GetSourceAccountID() string {
	if o == nil {
		return ""
	}
	return o.SourceAccountID
}

func (o *Payment) GetStatus() PaymentStatus {
	if o == nil {
		return PaymentStatus("")
	}
	return o.Status
}

func (o *Payment) GetType() PaymentType {
	if o == nil {
		return PaymentType("")
	}
	return o.Type
}
