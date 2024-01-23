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

type TransferReversalID struct {
	Reference   string
	ConnectorID ConnectorID
}

func (tid TransferReversalID) String() string {
	data, err := canonicaljson.Marshal(tid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func TransferReversalIDFromString(value string) (TransferReversalID, error) {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return TransferReversalID{}, err
	}
	ret := TransferReversalID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return TransferReversalID{}, err
	}

	return ret, nil
}

func MustTransferReversalIDFromString(value string) TransferReversalID {
	id, err := TransferReversalIDFromString(value)
	if err != nil {
		panic(err)
	}
	return id
}

func (tid TransferReversalID) Value() (driver.Value, error) {
	return tid.String(), nil
}

func (tid *TransferReversalID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := TransferReversalIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse paymentid %s: %v", v, err)
			}

			*tid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan paymentid: %v", value)
}

type TransferReversal struct {
	bun.BaseModel `bun:"transfers.transfer_reversal"`

	ID                   TransferReversalID `bun:",pk"`
	TransferInitiationID TransferInitiationID

	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string

	ConnectorID ConnectorID

	Amount *big.Int
	Asset  Asset

	Status TransferReversalStatus
	Error  string

	Metadata map[string]string
}
