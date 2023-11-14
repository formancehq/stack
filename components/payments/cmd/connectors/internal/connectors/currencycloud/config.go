package currencycloud

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"
)

type Config struct {
	Name          string              `json:"name" bson:"name"`
	LoginID       string              `json:"loginID" bson:"loginID"`
	APIKey        string              `json:"apiKey" bson:"apiKey"`
	Endpoint      string              `json:"endpoint" bson:"endpoint"`
	PollingPeriod connectors.Duration `json:"pollingPeriod" bson:"pollingPeriod"`
}

// String obfuscates sensitive fields and returns a string representation of the config.
// This is used for logging.
func (c Config) String() string {
	return fmt.Sprintf("loginID=%s, endpoint=%s, pollingPeriod=%s, apiKey=****", c.LoginID, c.Endpoint, c.PollingPeriod.String())
}

func (c Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}

	if c.LoginID == "" {
		return ErrMissingLoginID
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
	cfg.AddParameter("loginID", configtemplate.TypeString, true)
	cfg.AddParameter("apiKey", configtemplate.TypeString, true)
	cfg.AddParameter("endpoint", configtemplate.TypeString, false)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)

	return Name.String(), cfg
}
