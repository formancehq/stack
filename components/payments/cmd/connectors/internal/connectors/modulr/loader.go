package modulr

import (
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/gorilla/mux"
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
	if cfg.PollingPeriod.Duration == 0 {
		cfg.PollingPeriod.Duration = defaultPollingPeriod
	}

	if cfg.Name == "" {
		cfg.Name = name.String()
	}

	if cfg.PageSize == 0 {
		cfg.PageSize = defaultPageSize
	}

	return cfg
}

func (l *Loader) Router(_ *storage.Storage) *mux.Router {
	// Webhooks are not implemented yet
	return nil
}

// NewLoader creates a new loader.
func NewLoader() *Loader {
	return &Loader{}
}
