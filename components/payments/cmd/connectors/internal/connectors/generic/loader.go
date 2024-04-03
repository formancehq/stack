package generic

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

	return cfg
}

func (l *Loader) Router(store *storage.Storage) *mux.Router {
	return nil
}

// NewLoader creates a new loader.
func NewLoader() *Loader {
	return &Loader{}
}
