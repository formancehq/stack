package api

import (
	"net/http"

	"github.com/formancehq/payments/cmd/api/internal/storage"
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

	authGroup.Path("/payments").Methods(http.MethodGet).Handler(listPaymentsHandler(store))
	authGroup.Path("/payments/{paymentID}").Methods(http.MethodGet).Handler(readPaymentHandler(store))
	authGroup.Path("/payments/{paymentID}/metadata").Methods(http.MethodPatch).Handler(updateMetadataHandler(store))

	authGroup.Path("/accounts").Methods(http.MethodGet).Handler(listAccountsHandler(store))
	authGroup.Path("/accounts/{accountID}").Methods(http.MethodGet).Handler(readAccountHandler(store))
	authGroup.Path("/accounts/{accountID}/balances").Methods(http.MethodGet).Handler(listBalancesForAccount(store))

	authGroup.Path("/bank-accounts").Methods(http.MethodGet).Handler(listBankAccountsHandler(store))
	authGroup.Path("/bank-accounts/{bankAccountID}").Methods(http.MethodGet).Handler(readBankAccountHandler(store))

	authGroup.Path("/transfer-initiations").Methods(http.MethodGet).Handler(listTransferInitiationsHandler(store))
	authGroup.Path("/transfer-initiations/{transferID}").Methods(http.MethodGet).Handler(readTransferInitiationHandler(store))

	return rootMux, nil
}
