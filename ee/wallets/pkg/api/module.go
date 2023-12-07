package api

import (
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	sharedhealth "github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func Module(serviceInfo sharedapi.ServiceInfo, listen string) fx.Option {
	return fx.Module(
		"api",
		fx.Provide(NewRouter),
		fx.Supply(serviceInfo),
		sharedhealth.Module(),
		fx.Invoke(func(lc fx.Lifecycle, router *chi.Mux) {
			lc.Append(httpserver.NewHook(router, httpserver.WithAddress(listen)))
		}),
	)
}
