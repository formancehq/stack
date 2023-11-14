package api

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/connectors/internal/integration"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func httpRouter(
	logger logging.Logger,
	store *storage.Storage,
	serviceInfo api.ServiceInfo,
	publisher message.Publisher,
	connectorHandlers []connectorHandler,
) (*mux.Router, error) {
	rootMux := mux.NewRouter()

	// We have to keep this recovery handler here to ensure that the health
	// endpoint is not panicking
	rootMux.Use(recoveryHandler(httpRecoveryFunc))
	rootMux.Use(httpCorsHandler())
	rootMux.Use(httpServeFunc)
	rootMux.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r.WithContext(logging.ContextWithLogger(r.Context(), logger)))
		})
	})

	rootMux.Path("/_health").Handler(healthHandler(store))

	subRouter := rootMux.NewRoute().Subrouter()
	if viper.GetBool(otelTracesFlag) {
		subRouter.Use(otelmux.Middleware(serviceName))
		// Add a second recovery handler to ensure that the otel middleware
		// is catching the error in the trace
		rootMux.Use(recoveryHandler(httpRecoveryFunc))
	}
	subRouter.Path("/_live").Handler(liveHandler())
	subRouter.Path("/_info").Handler(api.InfoHandler(serviceInfo))

	authGroup := subRouter.Name("authenticated").Subrouter()

	authGroup.Path("/bank-accounts").Methods(http.MethodPost).Handler(createBankAccountHandler(store, publisher))

	paymentsHandlers := make(map[models.ConnectorProvider]paymentHandler)
	for _, h := range connectorHandlers {
		paymentsHandlers[h.Provider] = h.initiatePayment
	}
	authGroup.Path("/transfer-initiations").Methods(http.MethodPost).Handler(createTransferInitiationHandler(store, publisher, paymentsHandlers))
	authGroup.Path("/transfer-initiations/{transferID}/status").Methods(http.MethodPost).Handler(updateTransferInitiationStatusHandler(store, publisher, paymentsHandlers))
	authGroup.Path("/transfer-initiations/{transferID}/retry").Methods(http.MethodPost).Handler(retryTransferInitiationHandler(store, publisher, paymentsHandlers))
	authGroup.Path("/transfer-initiations/{transferID}").Methods(http.MethodDelete).Handler(deleteTransferInitiationHandler(store, publisher))

	authGroup.HandleFunc("/connectors", readConnectorsHandler(store))

	connectorGroup := authGroup.PathPrefix("/connectors").Subrouter()
	connectorGroup.Path("/configs").Handler(connectorConfigsHandler())

	for _, h := range connectorHandlers {
		connectorGroup.PathPrefix("/" + h.Provider.String()).Handler(
			http.StripPrefix("/connectors", h.Handler))

		connectorGroup.PathPrefix("/" + h.Provider.StringLower()).Handler(
			http.StripPrefix("/connectors", h.Handler))
	}

	return rootMux, nil
}

func connectorRouter[Config models.ConnectorConfigObject](
	provider models.ConnectorProvider,
	manager *integration.ConnectorsManager[Config],
) *mux.Router {
	r := mux.NewRouter()

	addRoute(r, provider, "", http.MethodPost, install(manager))
	addRoute(r, provider, "/{connectorID}", http.MethodDelete, uninstall(manager, V1))
	addRoute(r, provider, "/{connectorID}/config", http.MethodGet, readConfig(manager, V1))
	addRoute(r, provider, "/{connectorID}/reset", http.MethodPost, reset(manager, V1))
	addRoute(r, provider, "/{connectorID}/tasks", http.MethodGet, listTasks(manager, V1))
	addRoute(r, provider, "/{connectorID}/tasks/{taskID}", http.MethodGet, readTask(manager, V1))

	// Deprecated routes
	addRoute(r, provider, "", http.MethodDelete, uninstall(manager, V0))
	addRoute(r, provider, "/config", http.MethodGet, readConfig(manager, V0))
	addRoute(r, provider, "/reset", http.MethodPost, reset(manager, V0))
	addRoute(r, provider, "/tasks", http.MethodGet, listTasks(manager, V0))
	addRoute(r, provider, "/tasks/{taskID}", http.MethodGet, readTask(manager, V0))

	return r
}

func addRoute(r *mux.Router, provider models.ConnectorProvider, path, method string, handler http.Handler) {
	r.Path("/" + provider.String() + path).Methods(method).Handler(handler)
	r.Path("/" + provider.StringLower() + path).Methods(method).Handler(handler)
}
