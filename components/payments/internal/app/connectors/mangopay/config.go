package mangopay

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors"
	"github.com/formancehq/payments/internal/app/connectors/configtemplate"
)

type Config struct {
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

	return nil
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("clientID", configtemplate.TypeString, true)
	cfg.AddParameter("apiKey", configtemplate.TypeString, true)
	cfg.AddParameter("endpoint", configtemplate.TypeString, true)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)

	return Name.String(), cfg
}
