package api

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/dummypay"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	otelTracesFlag = "otel-traces"
	serviceName    = "Payments"

	ErrUniqueReference      = "CONFLICT"
	ErrNotFound             = "NOT_FOUND"
	ErrInvalidID            = "INVALID_ID"
	ErrMissingOrInvalidBody = "MISSING_OR_INVALID_BODY"
	ErrValidation           = "VALIDATION"
)

func HTTPModule(serviceInfo api.ServiceInfo, bind string) fx.Option {
	return fx.Options(
		fx.Invoke(func(m *mux.Router, lc fx.Lifecycle) {
			lc.Append(httpserver.NewHook(m, httpserver.WithAddress(bind)))
		}),
		fx.Supply(serviceInfo),
		fx.Provide(fx.Annotate(httpRouter, fx.ParamTags(``, ``, ``, ``, `group:"connectorHandlers"`))),
		addConnector[dummypay.Config](dummypay.NewLoader()),
		addConnector[modulr.Config](modulr.NewLoader()),
		addConnector[stripe.Config](stripe.NewLoader()),
		addConnector[wise.Config](wise.NewLoader()),
		addConnector[currencycloud.Config](currencycloud.NewLoader()),
		addConnector[bankingcircle.Config](bankingcircle.NewLoader()),
		addConnector[mangopay.Config](mangopay.NewLoader()),
		addConnector[moneycorp.Config](moneycorp.NewLoader()),
	)
}

func httpRecoveryFunc(ctx context.Context, e interface{}) {
	if viper.GetBool(otelTracesFlag) {
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

func handleStorageErrors(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, storage.ErrDuplicateKeyValue) {
		api.BadRequest(w, ErrUniqueReference, err)
	} else if errors.Is(err, storage.ErrNotFound) {
		api.NotFound(w)
	} else {
		api.InternalServerError(w, r, err)
	}
}

func pageSizeQueryParam(r *http.Request) (int, error) {
	if value := r.URL.Query().Get("pageSize"); value != "" {
		ret, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return 0, err
		}

		return int(ret), nil
	}

	return 0, nil
}
