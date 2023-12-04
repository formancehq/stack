package models

import (
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/gibson042/canonicaljson-go"
	"github.com/uptrace/bun"
)

type TransferInitiationID struct {
	Reference   string
	ConnectorID ConnectorID
}

func (tid TransferInitiationID) String() string {
	data, err := canonicaljson.Marshal(tid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func TransferInitiationIDFromString(value string) (TransferInitiationID, error) {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return TransferInitiationID{}, err
	}
	ret := TransferInitiationID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return TransferInitiationID{}, err
	}

	return ret, nil
}

func MustTransferInitiationIDFromString(value string) TransferInitiationID {
	id, err := TransferInitiationIDFromString(value)
	if err != nil {
		panic(err)
	}
	return id
}

func (tid TransferInitiationID) Value() (driver.Value, error) {
	return tid.String(), nil
}

func (tid *TransferInitiationID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := TransferInitiationIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse paymentid %s: %v", v, err)
			}

			*tid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan paymentid: %v", value)
}

type TransferInitiationStatus int

const (
	TransferInitiationStatusWaitingForValidation TransferInitiationStatus = iota
	TransferInitiationStatusProcessing
	TransferInitiationStatusProcessed
	TransferInitiationStatusFailed
	TransferInitiationStatusRejected
	TransferInitiationStatusValidated
)

func (s TransferInitiationStatus) String() string {
	return [...]string{
		"WAITING_FOR_VALIDATION",
		"PROCESSING",
		"PROCESSED",
		"FAILED",
		"REJECTED",
		"VALIDATED",
	}[s]
}

func TransferInitiationStatusFromString(s string) (TransferInitiationStatus, error) {
	switch s {
	case "WAITING_FOR_VALIDATION":
		return TransferInitiationStatusWaitingForValidation, nil
	case "PROCESSING":
		return TransferInitiationStatusProcessing, nil
	case "PROCESSED":
		return TransferInitiationStatusProcessed, nil
	case "FAILED":
		return TransferInitiationStatusFailed, nil
	case "REJECTED":
		return TransferInitiationStatusRejected, nil
	case "VALIDATED":
		return TransferInitiationStatusValidated, nil
	default:
		return TransferInitiationStatusWaitingForValidation, errors.New("invalid status")
	}
}

type TransferInitiationType int

const (
	TransferInitiationTypeTransfer TransferInitiationType = iota
	TransferInitiationTypePayout
)

func (t TransferInitiationType) String() string {
	return [...]string{
		"TRANSFER",
		"PAYOUT",
	}[t]
}

func TransferInitiationTypeFromString(s string) (TransferInitiationType, error) {
	switch s {
	case "TRANSFER":
		return TransferInitiationTypeTransfer, nil
	case "PAYOUT":
		return TransferInitiationTypePayout, nil
	default:
		return TransferInitiationTypeTransfer, errors.New("invalid type")
	}
}

func MustTransferInitiationTypeFromString(s string) TransferInitiationType {
	t, err := TransferInitiationTypeFromString(s)
	if err != nil {
		panic(err)
	}
	return t
}

type TransferInitiation struct {
	bun.BaseModel `bun:"transfers.transfer_initiation"`

	// Filled when created in DB
	ID TransferInitiationID `bun:",pk,nullzero"`

	CreatedAt   time.Time `bun:",nullzero"`
	ScheduledAt time.Time `bun:",nullzero"`
	UpdatedAt   time.Time `bun:",nullzero"`
	Description string

	Type TransferInitiationType

	SourceAccountID      *AccountID
	DestinationAccountID AccountID
	Provider             ConnectorProvider
	ConnectorID          ConnectorID

	Amount *big.Int `bun:"type:numeric"`
	Asset  Asset

	Attempts int

	// We still have to keep the status and error here in case the transfer
	// failed before creating a payment (i.e. the transfer was rejected)
	Status TransferInitiationStatus
	Error  string

	SourceAccount      *Account `bun:"-"`
	DestinationAccount *Account `bun:"-"`

	RelatedPayments []*TransferInitiationPayments `bun:"-"`
}

type TransferInitiationPayments struct {
	bun.BaseModel `bun:"transfers.transfer_initiation_payments"`

	TransferInitiationID TransferInitiationID `bun:",pk"`
	PaymentID            PaymentID            `bun:",pk"`

	CreatedAt time.Time `bun:",nullzero"`
	Status    TransferInitiationStatus
	Error     string
}
