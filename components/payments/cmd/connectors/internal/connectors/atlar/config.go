package atlar

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Config struct {
	Name          string              `json:"name" yaml:"name" bson:"name"`
	PollingPeriod connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
	BaseUrl       string              `json:"baseUrl" yaml:"baseUrl" bson:"baseUrl"`
	AccessKey     string              `json:"accessKey" yaml:"accessKey" bson:"accessKey"`
	Secret        string              `json:"secret" yaml:"secret" bson:"secret"`
	ApiConfig     `bson:",inline"`
}

// String obfuscates sensitive fields and returns a string representation of the config.
// This is used for logging.
func (c Config) String() string {
	return fmt.Sprintf("baseUrl=%s, pollingPeriod=%s, pageSize=%d, accessKey=%s, secret=****",
		c.BaseUrl, c.PollingPeriod, c.PageSize, c.AccessKey)
}

func (c Config) Validate() error {
	if c.AccessKey == "" {
		return errors.New("missing api access key")
	}

	if c.Secret == "" {
		return errors.New("missing api secret")
	}

	return nil
}

func (c Config) ConnectorName() string {
	return c.Name
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

type ApiConfig struct {
	PageSize uint64 `json:"pageSize" yaml:"pageSize" bson:"pageSize"`
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("baseUrl", configtemplate.TypeString, false)
	cfg.AddParameter("accessKey", configtemplate.TypeString, true)
	cfg.AddParameter("secret", configtemplate.TypeString, true)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)
	cfg.AddParameter("pageSize", configtemplate.TypeDurationUnsignedInteger, false)

	return Name.String(), cfg
}
