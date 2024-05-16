package api

import (
	"context"

	"github.com/formancehq/reconciliation/internal/api/client"
	"github.com/formancehq/reconciliation/internal/api/v1/service"
	"github.com/formancehq/reconciliation/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
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

func Module(serviceInfo api.ServiceInfo, bind string) fx.Option {
	return fx.Options(
		healthCheckModule(),
		client.Module(),
		fx.Supply(serviceInfo),
		fx.Invoke(func(m *chi.Mux, lc fx.Lifecycle) {
			lc.Append(httpserver.NewHook(m, httpserver.WithAddress(bind)))
		}),
		fx.Provide(func(store *storage.Storage) service.Store {
			return store
		}),
		fx.Provide(newRouter),
	)
}
