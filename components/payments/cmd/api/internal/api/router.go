package api

import (
	"net/http"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func httpRouter(
	b backend.Backend,
	logger logging.Logger,
	serviceInfo api.ServiceInfo,
	a auth.Auth,
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
	if viper.GetBool(otelTracesFlag) {
		subRouter.Use(otelmux.Middleware(serviceName))
		// Add a second recovery handler to ensure that the otel middleware
		// is catching the error in the trace
		rootMux.Use(recoveryHandler(httpRecoveryFunc))
	}
	subRouter.Path("/_live").Handler(liveHandler())
	subRouter.Path("/_info").Handler(api.InfoHandler(serviceInfo))

	authGroup := subRouter.Name("authenticated").Subrouter()
	authGroup.Use(auth.Middleware(a))

	authGroup.Path("/payments").Methods(http.MethodPost).Handler(createPaymentHandler(b))
	authGroup.Path("/payments").Methods(http.MethodGet).Handler(listPaymentsHandler(b))
	authGroup.Path("/payments/{paymentID}").Methods(http.MethodGet).Handler(readPaymentHandler(b))
	authGroup.Path("/payments/{paymentID}/metadata").Methods(http.MethodPatch).Handler(updateMetadataHandler(b))

	authGroup.Path("/accounts").Methods(http.MethodGet).Handler(listAccountsHandler(b))
	authGroup.Path("/accounts/{accountID}").Methods(http.MethodGet).Handler(readAccountHandler(b))
	authGroup.Path("/accounts/{accountID}/balances").Methods(http.MethodGet).Handler(listBalancesForAccount(b))

	authGroup.Path("/bank-accounts").Methods(http.MethodGet).Handler(listBankAccountsHandler(b))
	authGroup.Path("/bank-accounts/{bankAccountID}").Methods(http.MethodGet).Handler(readBankAccountHandler(b))

	authGroup.Path("/transfer-initiations").Methods(http.MethodGet).Handler(listTransferInitiationsHandler(b))
	authGroup.Path("/transfer-initiations/{transferID}").Methods(http.MethodGet).Handler(readTransferInitiationHandler(b))

	authGroup.Path("/pools").Methods(http.MethodPost).Handler(createPoolHandler(b))
	authGroup.Path("/pools").Methods(http.MethodGet).Handler(listPoolHandler(b))
	authGroup.Path("/pools/{poolID}").Methods(http.MethodGet).Handler(getPoolHandler(b))
	authGroup.Path("/pools/{poolID}").Methods(http.MethodDelete).Handler(deletePoolHandler(b))
	authGroup.Path("/pools/{poolID}/accounts").Methods(http.MethodPost).Handler(addAccountToPoolHandler(b))
	authGroup.Path("/pools/{poolID}/accounts/{accountID}").Methods(http.MethodDelete).Handler(removeAccountFromPoolHandler(b))
	authGroup.Path("/pools/{poolID}/balances").Methods(http.MethodGet).Handler(getPoolBalances(b))

	return rootMux
}
