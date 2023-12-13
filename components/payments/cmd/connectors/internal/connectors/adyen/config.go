package adyen

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"
)

type Config struct {
	Name               string              `json:"name" yaml:"name" bson:"name"`
	APIKey             string              `json:"apiKey" yaml:"apiKey" bson:"apiKey"`
	HMACKey            string              `json:"hmacKey" yaml:"hmacKey" bson:"hmacKey"`
	LiveEndpointPrefix string              `json:"liveEndpointPrefix" yaml:"liveEndpointPrefix" bson:"liveEndpointPrefix"`
	PollingPeriod      connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
}

func (c Config) String() string {
	return fmt.Sprintf("liveEndpointPrefix=%s, apiKey=****, hmacKey=****", c.LiveEndpointPrefix)
}

func (c Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}

	if c.Name == "" {
		return ErrMissingName
	}

	if c.HMACKey == "" {
		return ErrMissingHMACKey
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
	cfg.AddParameter("apiKey", configtemplate.TypeString, true)
	cfg.AddParameter("hmacKey", configtemplate.TypeString, true)
	cfg.AddParameter("liveEndpointPrefix", configtemplate.TypeString, false)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)

	return Name.String(), cfg
}
