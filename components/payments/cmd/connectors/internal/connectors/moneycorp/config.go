package moneycorp

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"
)

const (
	pageSize = 100
)

type Config struct {
	Name          string              `json:"name" yaml:"name" bson:"name"`
	ClientID      string              `json:"clientID" yaml:"clientID" bson:"clientID"`
	APIKey        string              `json:"apiKey" yaml:"apiKey" bson:"apiKey"`
	Endpoint      string              `json:"endpoint" yaml:"endpoint" bson:"endpoint"`
	PollingPeriod connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
}

func (c Config) String() string {
	return fmt.Sprintf("clientID=%s, apiKey=****", c.ClientID)
}

func (c Config) Validate() error {
	if c.ClientID == "" {
		return ErrMissingClientID
	}

	if c.APIKey == "" {
		return ErrMissingAPIKey
	}

	if c.Endpoint == "" {
		return ErrMissingEndpoint
	}

	if c.Name == "" {
		return ErrMissingName
	}

	return nil
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c Config) ConnectorName() string {
	return c.Name
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("name", configtemplate.TypeString, true)
	cfg.AddParameter("clientID", configtemplate.TypeString, true)
	cfg.AddParameter("apiKey", configtemplate.TypeString, true)
	cfg.AddParameter("endpoint", configtemplate.TypeString, true)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)

	return Name.String(), cfg
}
