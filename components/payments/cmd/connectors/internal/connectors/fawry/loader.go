package fawry

import (
	"net/http"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/gorilla/mux"
)

type Loader struct{}

func (l *Loader) Name() models.ConnectorProvider {
	return "fawry"
}

func (l *Loader) AllowTasks() int {
	return 10
}

func (l *Loader) ApplyDefaults(cfg Config) Config {
	return cfg
}

// Router returns the router for the connector, which we'll use to handle webhooks.
func (l *Loader) Router(store *storage.Storage) *mux.Router {
	r := mux.NewRouter()

	w := NewWebhook()

	r.Path("/").Methods(http.MethodPost).HandlerFunc(w.Handle())

	return r
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(logger logging.Logger, config Config) Connector {
	return *newConnector()
}
