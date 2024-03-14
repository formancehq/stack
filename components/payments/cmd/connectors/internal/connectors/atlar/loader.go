package atlar

import (
	"net/url"

	"github.com/formancehq/payments/internal/models"
	"github.com/gorilla/mux"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

type Loader struct{}

const allowedTasks = 50

func (l *Loader) AllowTasks() int {
	return allowedTasks
}

func (l *Loader) Name() models.ConnectorProvider {
	return name
}

func (l *Loader) Load(logger logging.Logger, config Config) connectors.Connector {
	return newConnector(logger, config)
}

func (l *Loader) ApplyDefaults(cfg Config) Config {
	emptyUrl := url.URL{}
	if cfg.BaseUrl == emptyUrl {
		//"https://api.atlar.com"
		cfg.BaseUrl = defaultURLValue
	}

	if cfg.PageSize == 0 {
		cfg.PageSize = defaultPageSize
	}

	if cfg.PollingPeriod.Duration == 0 {
		cfg.PollingPeriod = connectors.Duration{Duration: defaultPollingPeriod}
	}

	if cfg.TransferInitiationStatusPollingPeriod.Duration == 0 {
		cfg.TransferInitiationStatusPollingPeriod = connectors.Duration{Duration: defaultPollingPeriod}
	}

	if cfg.Name == "" {
		cfg.Name = name.String()
	}

	return cfg
}

func (l *Loader) Router() *mux.Router {
	return nil
}

func NewLoader() *Loader {
	return &Loader{}
}
