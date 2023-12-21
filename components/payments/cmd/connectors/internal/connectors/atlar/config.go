package atlar

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Config struct {
	Name                                  string              `json:"name" yaml:"name" bson:"name"`
	PollingPeriod                         connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
	TransferInitiationStatusPollingPeriod connectors.Duration `json:"transferInitiationStatusPollingPeriod" yaml:"transferInitiationStatusPollingPeriod" bson:"transferInitiationStatusPollingPeriod"`
	BaseUrl                               url.URL             `json:"baseUrl" yaml:"baseUrl" bson:"baseUrl"`
	AccessKey                             string              `json:"accessKey" yaml:"accessKey" bson:"accessKey"`
	Secret                                string              `json:"secret" yaml:"secret" bson:"secret"`
	ApiConfig                             `bson:",inline"`
}

// String obfuscates sensitive fields and returns a string representation of the config.
// This is used for logging.
func (c Config) String() string {
	return fmt.Sprintf("baseUrl=%s, pollingPeriod=%s, transferInitiationStatusPollingPeriod=%s, pageSize=%d, accessKey=%s, secret=****",
		c.BaseUrl.String(), c.PollingPeriod, c.TransferInitiationStatusPollingPeriod, c.PageSize, c.AccessKey)
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
	type CopyType Config

	basicConfig := struct {
		BaseUrl string `json:"baseUrl"`
		CopyType
	}{
		BaseUrl:  c.BaseUrl.String(),
		CopyType: (CopyType)(c),
	}

	return json.Marshal(basicConfig)
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type CopyType Config

	tmp := struct {
		BaseUrl string `json:"baseUrl"`
		*CopyType
	}{
		CopyType: (*CopyType)(c),
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	baseUrl, err := url.Parse(tmp.BaseUrl)
	if err != nil {
		return err
	}
	c.BaseUrl = *baseUrl

	return nil
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
	cfg.AddParameter("transferInitiationStatusPollingPeriod", configtemplate.TypeDurationNs, false)
	cfg.AddParameter("pageSize", configtemplate.TypeDurationUnsignedInteger, false)

	return Name.String(), cfg
}
