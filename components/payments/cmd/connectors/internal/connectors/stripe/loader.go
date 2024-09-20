package stripe

import (
	"github.com/formancehq/payments/internal/models"
	"github.com/gorilla/mux"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
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
	if cfg.PageSize == 0 {
		cfg.PageSize = defaultPageSize
	}

	if cfg.PollingPeriod.Duration == 0 {
		cfg.PollingPeriod = connectors.Duration{Duration: defaultPollingPeriod}
	}

	if cfg.Name == "" {
		cfg.Name = name.String()
	}

	return cfg
}

func (l *Loader) Router(_ *storage.Storage) *mux.Router {
	// Webhooks are not implemented yet
	return nil
}

func NewLoader() *Loader {
	return &Loader{}
}
