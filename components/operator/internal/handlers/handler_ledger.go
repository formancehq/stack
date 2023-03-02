package handlers

import (
	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

func init() {
	modules.Register("ledger", modules.Module{
		Postgres: func(ctx modules.Context) v1beta3.PostgresConfig {
			return ctx.Configuration.Spec.Services.Ledger.Postgres
		},
		Services: func(ctx modules.Context) modules.Services {
			return modules.Services{{
				Port:                    8080,
				InjectPostgresVariables: true,
				HasVersionEndpoint:      true,
				Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
					env := modules.NewEnv().Append(
						modules.Env("BIND", "0.0.0.0:8080"),
						modules.Env("STORAGE_DRIVER", "postgres"),
						modules.Env("PUBLISHER_TOPIC_MAPPING", "*:"+resolveContext.Stack.GetServiceName("ledger")),
					).Append(modules.BrokerEnvVars(resolveContext.Configuration.Spec.Broker, "ledger")...)

					return modules.Container{
						Image: modules.GetImage("ledger", resolveContext.Versions.Spec.Ledger),
						Env: env.Append(
							modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
						),
					}
				},
			}}
		},
	})
}
