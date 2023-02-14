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
				EnvPrefix:               "NUMARY_",
				Port:                    8080,
				InjectPostgresVariables: true,
				HasVersionEndpoint:      true,
				Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
					env := modules.NewEnv().Append(
						modules.Env("NUMARY_SERVER_HTTP_BIND_ADDRESS", "0.0.0.0:8080"),
						modules.Env("NUMARY_STORAGE_DRIVER", "postgres"),
						modules.Env("NUMARY_PUBLISHER_TOPIC_MAPPING", "*:"+resolveContext.Stack.GetServiceName("ledger")),
					).Append(modules.BrokerEnvVarsWithPrefix(resolveContext.Configuration.Spec.Broker, "ledger", "NUMARY_")...)

					return modules.Container{
						Image: modules.GetImage("ledger", resolveContext.Versions.Spec.Ledger),
						Env: env.Append(
							modules.Env("NUMARY_STORAGE_POSTGRES_CONN_STRING", "$(NUMARY_POSTGRES_URI)"),
						),
					}
				},
			}}
		},
	})
}
