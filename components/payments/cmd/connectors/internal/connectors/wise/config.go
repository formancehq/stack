package wise

import (
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"
)

const (
	defaultPollingPeriod = 2 * time.Minute
	pageSize             = 100
)

type Config struct {
	Name          string              `json:"name" yaml:"name" bson:"name"`
	APIKey        string              `json:"apiKey" yaml:"apiKey" bson:"apiKey"`
	PollingPeriod connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
}

// String obfuscates sensitive fields and returns a string representation of the config.
// This is used for logging.
func (c Config) String() string {
	return "apiKey=***"
}

func (c Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
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

	cfg.AddParameter("name", configtemplate.TypeString, name.String(), false)
	cfg.AddParameter("apiKey", configtemplate.TypeString, "", true)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, defaultPollingPeriod.String(), false)

	return name.String(), cfg
}
