package modulr

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors"
	"github.com/formancehq/payments/internal/app/connectors/configtemplate"
)

type Config struct {
	APIKey        string              `json:"apiKey" bson:"apiKey"`
	APISecret     string              `json:"apiSecret" bson:"apiSecret"`
	Endpoint      string              `json:"endpoint" bson:"endpoint"`
	PollingPeriod connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
}

// String obfuscates sensitive fields and returns a string representation of the config.
// This is used for logging.
func (c Config) String() string {
	return fmt.Sprintf("endpoint=%s, apiSecret=***, apiKey=****", c.Endpoint)
}

func (c Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}

	if c.APISecret == "" {
		return ErrMissingAPISecret
	}

	return nil
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("apiKey", configtemplate.TypeString, true)
	cfg.AddParameter("apiSecret", configtemplate.TypeString, true)
	cfg.AddParameter("endpoint", configtemplate.TypeString, false)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)

	return Name.String(), cfg
}
