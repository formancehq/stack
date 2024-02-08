package api

import (
	"net/http"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func httpRouter(
	logger logging.Logger,
	b backend.ServiceBackend,
	serviceInfo api.ServiceInfo,
	a auth.Auth,
	connectorHandlers []connectorHandler,
) *mux.Router {
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

	rootMux.Path("/_health").Handler(healthHandler(b))

	subRouter := rootMux.NewRoute().Subrouter()
	if viper.GetBool(otlptraces.OtelTracesFlag) {
		subRouter.Use(otelmux.Middleware(serviceName))
		// Add a second recovery handler to ensure that the otel middleware
		// is catching the error in the trace
		rootMux.Use(recoveryHandler(httpRecoveryFunc))
	}
	subRouter.Path("/_live").Handler(liveHandler())
	subRouter.Path("/_info").Handler(api.InfoHandler(serviceInfo))

	authGroup := subRouter.Name("authenticated").Subrouter()
	authGroup.Use(auth.Middleware(a))

	authGroup.Path("/bank-accounts").Methods(http.MethodPost).Handler(createBankAccountHandler(b))
	authGroup.Path("/bank-accounts/{bankAccountID}/forward").Methods(http.MethodPost).Handler(forwardBankAccountToConnector(b))
	authGroup.Path("/bank-accounts/{bankAccountID}/metadata").Methods(http.MethodPatch).Handler(updateBankAccountMetadataHandler(b))

	authGroup.Path("/transfer-initiations").Methods(http.MethodPost).Handler(createTransferInitiationHandler(b))
	authGroup.Path("/transfer-initiations/{transferID}/status").Methods(http.MethodPost).Handler(updateTransferInitiationStatusHandler(b))
	authGroup.Path("/transfer-initiations/{transferID}/retry").Methods(http.MethodPost).Handler(retryTransferInitiationHandler(b))
	authGroup.Path("/transfer-initiations/{transferID}/reverse").Methods(http.MethodPost).Handler(reverseTransferInitiationHandler(b))
	authGroup.Path("/transfer-initiations/{transferID}").Methods(http.MethodDelete).Handler(deleteTransferInitiationHandler(b))

	authGroup.HandleFunc("/connectors", readConnectorsHandler(b))

	connectorGroupAuthenticated := authGroup.PathPrefix("/connectors").Subrouter()
	connectorGroupAuthenticated.Path("/configs").Handler(connectorConfigsHandler())

	// Needed for webhooks
	connectorGroupUnauthenticated := subRouter.PathPrefix("/connectors").Subrouter()

	for _, h := range connectorHandlers {
		connectorGroupAuthenticated.PathPrefix("/" + h.Provider.String()).Handler(
			http.StripPrefix("/connectors", h.Handler))

		connectorGroupAuthenticated.PathPrefix("/" + h.Provider.StringLower()).Handler(
			http.StripPrefix("/connectors", h.Handler))

		if h.WebhookHandler != nil {
			connectorGroupUnauthenticated.PathPrefix("/webhooks/" + h.Provider.String()).Handler(
				http.StripPrefix("/connectors", h.WebhookHandler))

			connectorGroupUnauthenticated.PathPrefix("/webhooks/" + h.Provider.StringLower()).Handler(
				http.StripPrefix("/connectors", h.WebhookHandler))
		}
	}

	return rootMux
}

func connectorRouter[Config models.ConnectorConfigObject](
	provider models.ConnectorProvider,
	b backend.ManagerBackend[Config],
) *mux.Router {
	r := mux.NewRouter()

	addRoute(r, provider, "", http.MethodPost, install(b))
	addRoute(r, provider, "/{connectorID}", http.MethodDelete, uninstall(b, V1))
	addRoute(r, provider, "/{connectorID}/config", http.MethodGet, readConfig(b, V1))
	addRoute(r, provider, "/{connectorID}/config", http.MethodPost, updateConfig(b, V1))
	addRoute(r, provider, "/{connectorID}/reset", http.MethodPost, reset(b, V1))
	addRoute(r, provider, "/{connectorID}/tasks", http.MethodGet, listTasks(b, V1))
	addRoute(r, provider, "/{connectorID}/tasks/{taskID}", http.MethodGet, readTask(b, V1))

	// Deprecated routes
	addRoute(r, provider, "", http.MethodDelete, uninstall(b, V0))
	addRoute(r, provider, "/config", http.MethodGet, readConfig(b, V0))
	addRoute(r, provider, "/reset", http.MethodPost, reset(b, V0))
	addRoute(r, provider, "/tasks", http.MethodGet, listTasks(b, V0))
	addRoute(r, provider, "/tasks/{taskID}", http.MethodGet, readTask(b, V0))

	return r
}

func webhookConnectorRouter[Config models.ConnectorConfigObject](
	provider models.ConnectorProvider,
	connectorRouter *mux.Router,
	b backend.ManagerBackend[Config],
) *mux.Router {
	if connectorRouter == nil {
		return nil
	}

	r := mux.NewRouter()

	group := r.PathPrefix("/webhooks/" + provider.String() + "/{connectorID}").Subrouter()
	group.Use(webhooksMiddleware(b, V1))
	addWebhookRoute(group, connectorRouter)

	groupLower := r.PathPrefix("/webhooks/" + provider.StringLower() + "/{connectorID}").Subrouter()
	groupLower.Use(webhooksMiddleware(b, V1))
	addWebhookRoute(groupLower, connectorRouter)

	return r
}

func addRoute(r *mux.Router, provider models.ConnectorProvider, path, method string, handler http.Handler) {
	r.Path("/" + provider.String() + path).Methods(method).Handler(handler)
	r.Path("/" + provider.StringLower() + path).Methods(method).Handler(handler)
}

func addWebhookRoute(r *mux.Router, subRouter *mux.Router) {
	subRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		methods, err := route.GetMethods()
		if err != nil {
			return err
		}

		for _, method := range methods {
			r.Path(pathTemplate).Methods(method).Handler(route.GetHandler())
		}

		return nil
	})
}
