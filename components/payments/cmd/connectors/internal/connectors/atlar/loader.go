package atlar

import (
	"net/url"
	"time"

	"github.com/formancehq/payments/internal/models"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

type Loader struct{}

const allowedTasks = 50

func (l *Loader) AllowTasks() int {
	return allowedTasks
}

func (l *Loader) Name() models.ConnectorProvider {
	return Name
}

func (l *Loader) Load(logger logging.Logger, config Config) connectors.Connector {
	return newConnector(logger, config)
}

func (l *Loader) ApplyDefaults(cfg Config) Config {
	defaultUrl := url.URL{}
	if cfg.BaseUrl == defaultUrl {
		//"https://api.atlar.com"
		cfg.BaseUrl = url.URL{
			Scheme: "https",
			Host:   "api.atlar.com",
		}
	}

	if cfg.PageSize == 0 {
		cfg.PageSize = 25
	}

	if cfg.PollingPeriod.Duration == 0 {
		cfg.PollingPeriod = connectors.Duration{Duration: 2 * time.Minute}
	}

	if cfg.Name == "" {
		cfg.Name = Name.String()
	}

	return cfg
}

func NewLoader() *Loader {
	return &Loader{}
}
