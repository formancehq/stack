package currencycloud

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors"
	"github.com/formancehq/payments/internal/app/connectors/configtemplate"
)

type Config struct {
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

	return nil
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("loginID", configtemplate.TypeString, true)
	cfg.AddParameter("apiKey", configtemplate.TypeString, true)
	cfg.AddParameter("endpoint", configtemplate.TypeString, false)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)

	return Name.String(), cfg
}
