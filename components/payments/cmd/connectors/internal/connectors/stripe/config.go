package stripe

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

const (
	defaultPageSize      = 10
	defaultPollingPeriod = 2 * time.Minute
)

type Config struct {
	Name           string              `json:"name" yaml:"name" bson:"name"`
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

	if c.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

func (c Config) ConnectorName() string {
	return c.Name
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

type TimelineConfig struct {
	PageSize uint64 `json:"pageSize" yaml:"pageSize" bson:"pageSize"`
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("name", configtemplate.TypeString, name.String(), false)
	cfg.AddParameter("apiKey", configtemplate.TypeString, "", true)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, defaultPollingPeriod.String(), false)
	cfg.AddParameter("pageSize", configtemplate.TypeDurationUnsignedInteger, strconv.Itoa(defaultPageSize), false)

	return name.String(), cfg
}
