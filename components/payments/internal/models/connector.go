package models

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gibson042/canonicaljson-go"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Connector struct {
	bun.BaseModel `bun:"connectors.connector"`

	ID        ConnectorID `bun:",pk,nullzero"`
	Name      string
	CreatedAt time.Time `bun:",nullzero"`
	Provider  ConnectorProvider

	// EncryptedConfig is a PGP-encrypted JSON string.
	EncryptedConfig string `bun:"config"`

	// Config is a decrypted config. It is not stored in the database.
	Config json.RawMessage `bun:"decrypted_config,scanonly"`

	Tasks []*Task `bun:"rel:has-many,join:id=connector_id"`
}

type ConnectorID struct {
	Reference uuid.UUID
	Provider  ConnectorProvider
}

func (cid *ConnectorID) String() string {
	if cid == nil || cid.Reference == uuid.Nil {
		return ""
	}

	data, err := canonicaljson.Marshal(cid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func ConnectorIDFromString(value string) (ConnectorID, error) {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return ConnectorID{}, err
	}
	ret := ConnectorID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return ConnectorID{}, err
	}

	return ret, nil
}

func MustConnectorIDFromString(value string) ConnectorID {
	id, err := ConnectorIDFromString(value)
	if err != nil {
		panic(err)
	}
	return id
}

func (cid ConnectorID) Value() (driver.Value, error) {
	return cid.String(), nil
}

func (cid *ConnectorID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("connector id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := ConnectorIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse connector id %s: %v", v, err)
			}

			*cid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan connector id: %v", value)
}

func (c Connector) String() string {
	c.EncryptedConfig = "****"
	c.Config = nil

	var t any = c

	return fmt.Sprintf("%+v", t)
}

type ConnectorProvider string

const (
	ConnectorProviderBankingCircle ConnectorProvider = "BANKING-CIRCLE"
	ConnectorProviderCurrencyCloud ConnectorProvider = "CURRENCY-CLOUD"
	ConnectorProviderDummyPay      ConnectorProvider = "DUMMY-PAY"
	ConnectorProviderModulr        ConnectorProvider = "MODULR"
	ConnectorProviderStripe        ConnectorProvider = "STRIPE"
	ConnectorProviderWise          ConnectorProvider = "WISE"
	ConnectorProviderMangopay      ConnectorProvider = "MANGOPAY"
	ConnectorProviderMoneycorp     ConnectorProvider = "MONEYCORP"
)

func (p ConnectorProvider) String() string {
	return string(p)
}

func (p ConnectorProvider) StringLower() string {
	return strings.ToLower(string(p))
}

func ConnectorProviderFromString(s string) (ConnectorProvider, error) {
	switch s {
	case "BANKING-CIRCLE":
		return ConnectorProviderBankingCircle, nil
	case "CURRENCY-CLOUD":
		return ConnectorProviderCurrencyCloud, nil
	case "DUMMY-PAY":
		return ConnectorProviderDummyPay, nil
	case "MODULR":
		return ConnectorProviderModulr, nil
	case "STRIPE":
		return ConnectorProviderStripe, nil
	case "WISE":
		return ConnectorProviderWise, nil
	case "MANGOPAY":
		return ConnectorProviderMangopay, nil
	case "MONEYCORP":
		return ConnectorProviderMoneycorp, nil
	default:
		return "", errors.New("unknown connector provider")
	}
}

func MustConnectorProviderFromString(s string) ConnectorProvider {
	p, err := ConnectorProviderFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

func (c Connector) ParseConfig(to interface{}) error {
	if c.Config == nil {
		return nil
	}

	err := json.Unmarshal(c.Config, to)
	if err != nil {
		return fmt.Errorf("failed to parse config (%s): %w", string(c.Config), err)
	}

	return nil
}

type ConnectorConfigObject interface {
	ConnectorName() string
	Validate() error
	Marshal() ([]byte, error)
}

type EmptyConnectorConfig struct {
	Name string
}

func (cfg EmptyConnectorConfig) ConnectorName() string {
	return cfg.Name
}

func (cfg EmptyConnectorConfig) Validate() error {
	return nil
}

func (cfg EmptyConnectorConfig) Marshal() ([]byte, error) {
	return nil, nil
}
