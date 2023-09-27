package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Connector struct {
	bun.BaseModel `bun:"connectors.connector"`

	ID        uuid.UUID `bun:",pk,nullzero"`
	CreatedAt time.Time `bun:",nullzero"`
	Provider  ConnectorProvider
	Enabled   bool

	// EncryptedConfig is a PGP-encrypted JSON string.
	EncryptedConfig string `bun:"config"`

	// Config is a decrypted config. It is not stored in the database.
	Config json.RawMessage `bun:"decrypted_config,scanonly"`

	Tasks    []*Task    `bun:"rel:has-many,join:id=connector_id"`
	Payments []*Payment `bun:"rel:has-many,join:id=connector_id"`
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
	Validate() error
	Marshal() ([]byte, error)
}

type EmptyConnectorConfig struct{}

func (cfg EmptyConnectorConfig) Validate() error {
	return nil
}

func (cfg EmptyConnectorConfig) Marshal() ([]byte, error) {
	return nil, nil
}
