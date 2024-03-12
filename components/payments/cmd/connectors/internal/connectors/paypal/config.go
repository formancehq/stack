package paypal

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"
	"github.com/formancehq/payments/internal/models"
)

var (
	_ models.ConnectorConfigObject = Config{}
)

type Config struct {
	Name          string              `json:"name" yaml:"name" bson:"name"`
	ClientID      string              `json:"clientID" yaml:"clientID" bson:"clientID"`
	Secret        string              `json:"secret" yaml:"secret" bson:"secret"`
	Endpoint      string              `json:"endpoint" yaml:"endpoint" bson:"endpoint"`
	PollingPeriod connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
}

func (c Config) Validate() (err error) {
	switch {
	case c.Name == "":
		err = ErrMissingName
	case c.ClientID == "":
		err = ErrMissingClientID
	case c.Secret == "":
		err = ErrMissingSecret
	case c.Endpoint == "":
		err = ErrMissingEndpoint
	}
	return
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()
	cfg.AddParameter("name", configtemplate.TypeString, true)
	cfg.AddParameter("clientID", configtemplate.TypeString, true)
	cfg.AddParameter("secret", configtemplate.TypeString, true)
	cfg.AddParameter("endpoint", configtemplate.TypeString, true)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)

	return models.ConnectorProviderPaypal.String(), cfg
}

func (c Config) ConnectorName() string {
	return c.Name
}

func (c Config) String() string {
	return fmt.Sprintf("clientID=%s, secret=****", c.ClientID)
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(&c)
}
