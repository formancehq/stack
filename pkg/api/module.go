package api

import (
	"net/url"

	"github.com/formancehq/auth/pkg/api/routing"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"go.uber.org/fx"
)

func Module(addr string, baseUrl *url.URL) fx.Option {
	return fx.Options(
		sharedhealth.ProvideHealthCheck(delegatedOIDCServerAvailability),
		routing.Module(addr, baseUrl),
		fx.Invoke(
			fx.Annotate(addClientRoutes, fx.ParamTags(``, `name:"prefixedRouter"`)),
			fx.Annotate(addScopeRoutes, fx.ParamTags(``, `name:"prefixedRouter"`)),
			fx.Annotate(addUserRoutes, fx.ParamTags(``, `name:"prefixedRouter"`)),
		),
	)
}
