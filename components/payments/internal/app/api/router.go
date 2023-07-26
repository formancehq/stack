package api

import (
	"net/http"

	"github.com/formancehq/payments/internal/app/integration"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func httpRouter(logger logging.Logger, store *storage.Storage, serviceInfo api.ServiceInfo, connectorHandlers []connectorHandler) (*mux.Router, error) {
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

	authGroup.Path("/payments").Methods(http.MethodGet).Handler(listPaymentsHandler(store))
	authGroup.Path("/payments/{paymentID}").Methods(http.MethodGet).Handler(readPaymentHandler(store))
	authGroup.Path("/payments/{paymentID}/metadata").Methods(http.MethodPatch).Handler(updateMetadataHandler(store))

	authGroup.Path("/accounts").Methods(http.MethodGet).Handler(listAccountsHandler(store))
	authGroup.Path("/accounts/{accountID}/balances").Methods(http.MethodGet).Handler(listBalancesForAccount(store))

	authGroup.HandleFunc("/connectors", readConnectorsHandler(store))

	connectorGroup := authGroup.PathPrefix("/connectors").Subrouter()

	connectorGroup.Path("/configs").Handler(connectorConfigsHandler())

	// Deprecated
	// TODO: Remove this endpoint
	// Use /connectors/stripe/transfers instead
	connectorGroup.Path("/stripe/transfers").Methods(http.MethodPost).
		Handler(handleStripeTransfers(store))

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
	manager *integration.ConnectorManager[Config],
) *mux.Router {
	r := mux.NewRouter()

	addRoute(r, provider, "", http.MethodPost, install(manager))
	addRoute(r, provider, "", http.MethodDelete, uninstall(manager))
	addRoute(r, provider, "/config", http.MethodGet, readConfig(manager))
	addRoute(r, provider, "/reset", http.MethodPost, reset(manager))
	addRoute(r, provider, "/tasks", http.MethodGet, listTasks(manager))
	addRoute(r, provider, "/tasks/{taskID}", http.MethodGet, readTask(manager))
	addRoute(r, provider, "/transfers", http.MethodPost, initiateTransfer(manager))
	addRoute(r, provider, "/transfers", http.MethodGet, listTransfers(manager))

	return r
}

func addRoute(r *mux.Router, provider models.ConnectorProvider, path, method string, handler http.Handler) {
	r.Path("/" + provider.String() + path).Methods(method).Handler(handler)

	r.Path("/" + provider.StringLower() + path).Methods(method).Handler(handler)
}
