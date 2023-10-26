package stripe

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Config struct {
	PollingPeriod  connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
	APIKey         string              `json:"apiKey" yaml:"apiKey" bson:"apiKey"`
	TimelineConfig `bson:",inline"`
}

// String obfuscates sensitive fields and returns a string representation of the config.
// This is used for logging.
func (c Config) String() string {
	return fmt.Sprintf("pollingPeriod=%s, pageSize=%d, apiKey=****", c.PollingPeriod, c.PageSize)
}

func (c Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("missing api key")
	}

	return nil
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

type TimelineConfig struct {
	PageSize uint64 `json:"pageSize" yaml:"pageSize" bson:"pageSize"`
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("apiKey", configtemplate.TypeString, true)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)
	cfg.AddParameter("pageSize", configtemplate.TypeDurationUnsignedInteger, false)

	return Name.String(), cfg
}
