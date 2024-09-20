package dummypay

import (
	"time"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/gorilla/mux"
)

type Loader struct{}

// Name returns the name of the connector.
func (l *Loader) Name() models.ConnectorProvider {
	return name
}

// AllowTasks returns the amount of tasks that are allowed to be scheduled.
func (l *Loader) AllowTasks() int {
	return 10
}

const (
	// defaultFilePollingPeriod is the default period between file polling.
	defaultFilePollingPeriod = 10 * time.Second
)

// ApplyDefaults applies default values to the configuration.
func (l *Loader) ApplyDefaults(cfg Config) Config {
	if cfg.FilePollingPeriod.Duration == 0 {
		cfg.FilePollingPeriod = connectors.Duration{Duration: defaultFilePollingPeriod}
	}

	if cfg.Name == "" {
		cfg.Name = name.String()
	}

	return cfg
}

// Load returns the connector.
func (l *Loader) Load(logger logging.Logger, config Config) connectors.Connector {
	return newConnector(logger, config, newFS())
}

// Router returns the router.
func (l *Loader) Router(_ *storage.Storage) *mux.Router {
	// Webhooks are not implemented yet
	return nil
}

// NewLoader creates a new loader.
func NewLoader() *Loader {
	return &Loader{}
}
