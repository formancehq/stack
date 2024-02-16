package api

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/adyen"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/dummypay"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	serviceName = "Payments"

	ErrUniqueReference      = "CONFLICT"
	ErrNotFound             = "NOT_FOUND"
	ErrInvalidID            = "INVALID_ID"
	ErrMissingOrInvalidBody = "MISSING_OR_INVALID_BODY"
	ErrValidation           = "VALIDATION"
)

func HTTPModule(serviceInfo api.ServiceInfo, bind, stackURL string) fx.Option {
	return fx.Options(
		fx.Invoke(func(m *mux.Router, lc fx.Lifecycle) {
			lc.Append(httpserver.NewHook(m, httpserver.WithAddress(bind)))
		}),
		fx.Supply(serviceInfo),
		fx.Provide(fx.Annotate(connectorsHandlerMap, fx.ParamTags(`group:"connectorHandlers"`))),
		fx.Provide(func(store *storage.Storage) service.Store {
			return store
		}),
		fx.Provide(fx.Annotate(service.New, fx.As(new(backend.Service)))),
		fx.Provide(backend.NewDefaultBackend),
		fx.Provide(fx.Annotate(httpRouter, fx.ParamTags(``, ``, ``, ``, `group:"connectorHandlers"`))),
		fx.Provide(func() *messages.Messages {
			return messages.NewMessages(stackURL)
		}),
		addConnector[dummypay.Config](dummypay.NewLoader()),
		addConnector[modulr.Config](modulr.NewLoader()),
		addConnector[stripe.Config](stripe.NewLoader()),
		addConnector[wise.Config](wise.NewLoader()),
		addConnector[currencycloud.Config](currencycloud.NewLoader()),
		addConnector[bankingcircle.Config](bankingcircle.NewLoader()),
		addConnector[mangopay.Config](mangopay.NewLoader()),
		addConnector[moneycorp.Config](moneycorp.NewLoader()),
		addConnector[atlar.Config](atlar.NewLoader()),
		addConnector[adyen.Config](adyen.NewLoader()),
	)
}

func connectorsHandlerMap(connectorHandlers []connectorHandler) map[models.ConnectorProvider]*service.ConnectorHandlers {
	m := make(map[models.ConnectorProvider]*service.ConnectorHandlers)
	for _, h := range connectorHandlers {
		if handlers, ok := m[h.Provider]; ok {
			handlers.InitiatePaymentHandler = h.initiatePayment
			handlers.ReversePaymentHandler = h.reversePayment
			handlers.BankAccountHandler = h.createExternalBankAccount
		} else {
			m[h.Provider] = &service.ConnectorHandlers{
				InitiatePaymentHandler: h.initiatePayment,
				ReversePaymentHandler:  h.reversePayment,
				BankAccountHandler:     h.createExternalBankAccount,
			}
		}
	}
	return m
}

func httpRecoveryFunc(ctx context.Context, e interface{}) {
	if viper.GetBool(otlptraces.OtelTracesFlag) {
		otlp.RecordAsError(ctx, e)
	} else {
		logrus.Errorln(e)
		debug.PrintStack()
	}
}

func httpCorsHandler() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut},
		AllowCredentials: true,
	}).Handler
}

func httpServeFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

func handleServiceErrors(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, storage.ErrDuplicateKeyValue):
		api.BadRequest(w, ErrUniqueReference, err)
	case errors.Is(err, storage.ErrNotFound):
		api.NotFound(w, err)
	case errors.Is(err, service.ErrValidation):
		api.BadRequest(w, ErrValidation, err)
	case errors.Is(err, service.ErrInvalidID):
		api.BadRequest(w, ErrInvalidID, err)
	case errors.Is(err, service.ErrPublish):
		api.InternalServerError(w, r, err)
	default:
		api.InternalServerError(w, r, err)
	}
}

func handleConnectorsManagerErrors(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, storage.ErrDuplicateKeyValue):
		api.BadRequest(w, ErrUniqueReference, err)
	case errors.Is(err, storage.ErrNotFound):
		api.NotFound(w, err)
	case errors.Is(err, manager.ErrAlreadyInstalled):
		api.BadRequest(w, ErrValidation, err)
	case errors.Is(err, manager.ErrNotInstalled):
		api.BadRequest(w, ErrValidation, err)
	case errors.Is(err, manager.ErrConnectorNotFound):
		api.BadRequest(w, ErrValidation, err)
	case errors.Is(err, manager.ErrNotFound):
		api.NotFound(w, err)
	case errors.Is(err, manager.ErrValidation):
		api.BadRequest(w, ErrValidation, err)
	default:
		api.InternalServerError(w, r, err)
	}
}
