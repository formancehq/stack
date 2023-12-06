package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/formancehq/reconciliation/internal/api/backend"
	"github.com/formancehq/reconciliation/internal/api/service"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

const (
	ErrInvalidID            = "INVALID_ID"
	ErrMissingOrInvalidBody = "MISSING_OR_INVALID_BODY"
	ErrValidation           = "VALIDATION"
)

func healthCheckModule() fx.Option {
	return fx.Options(
		health.Module(),
		health.ProvideHealthCheck(func() health.NamedCheck {
			return health.NewNamedCheck("default", health.CheckFn(func(ctx context.Context) error {
				return nil
			}))
		}),
	)
}

func HTTPModule(serviceInfo api.ServiceInfo, bind string) fx.Option {
	return fx.Options(
		healthCheckModule(),
		fx.Invoke(func(m *chi.Mux, lc fx.Lifecycle) {
			lc.Append(httpserver.NewHook(m, httpserver.WithAddress(bind)))
		}),
		fx.Supply(serviceInfo),
		fx.Provide(fx.Annotate(service.NewService, fx.As(new(backend.Service)))),
		fx.Provide(backend.NewDefaultBackend),
		fx.Provide(newRouter),
	)
}

func handleServiceErrors(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		api.BadRequest(w, ErrValidation, err)
	default:
		api.InternalServerError(w, r, err)
	}
}
