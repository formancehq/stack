package api

import (
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		sharedhealth.Module(),
		sharedhealth.ProvideHealthCheck(delegatedOIDCServerAvailability),
		fx.Invoke(
			addClientRoutes,
			addScopeRoutes,
			addUserRoutes,
		),
	)
}
