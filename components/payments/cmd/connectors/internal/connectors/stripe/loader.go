package stripe

import (
	"time"

	"github.com/formancehq/payments/internal/models"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/integration"
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

func (l *Loader) Load(logger logging.Logger, config Config) integration.Connector {
	return newConnector(logger, config)
}

func (l *Loader) ApplyDefaults(cfg Config) Config {
	if cfg.PageSize == 0 {
		cfg.PageSize = 10
	}

	if cfg.PollingPeriod.Duration == 0 {
		cfg.PollingPeriod = connectors.Duration{Duration: 2 * time.Minute}
	}

	return cfg
}

func NewLoader() *Loader {
	return &Loader{}
}
