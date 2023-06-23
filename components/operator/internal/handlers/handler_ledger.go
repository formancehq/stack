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
		Versions: map[string]modules.Version{
			"v0.0.0": {
				Services: func(ctx modules.ModuleContext) modules.Services {
					return modules.Services{{
						EnvPrefix:               "NUMARY_",
						ListenEnvVar:            "SERVER_HTTP_BIND_ADDRESS",
						InjectPostgresVariables: true,
						HasVersionEndpoint:      true,
						ExposeHTTP:              true,
						NeedTopic:               true,
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
							env := modules.NewEnv().Append(
								modules.Env("STORAGE_DRIVER", "postgres"),
								modules.Env("PUBLISHER_TOPIC_MAPPING", "*:"+resolveContext.Stack.GetServiceName("ledger")),
							).Append(modules.BrokerEnvVarsWithPrefix(resolveContext.Configuration.Spec.Broker, "ledger")...)
							if resolveContext.Configuration.Spec.Services.Ledger.AllowPastTimestamps {
								env = env.Append(modules.Env("COMMIT_POLICY", "allow-past-timestamps"))
							}

							return modules.Container{
								Resources: modules.ResourceSizeSmall(),
								Image:     modules.GetImage("ledger", resolveContext.Versions.Spec.Ledger),
								Env: env.Append(
									modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(NUMARY_POSTGRES_URI)"),
								),
							}
						},
					}}
				},
			},
			"v2.0.0": {
				Services: func(ctx modules.ModuleContext) modules.Services {
					return modules.Services{{
						ListenEnvVar:            "BIND",
						InjectPostgresVariables: true,
						HasVersionEndpoint:      true,
						ExposeHTTP:              true,
						NeedTopic:               true,
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
							env := modules.NewEnv().Append(
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
			},
		},
	})
}
