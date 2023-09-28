package handlers

import (
	"strconv"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

type ledgerModule struct{}

func (l ledgerModule) Postgres(ctx modules.Context) v1beta3.PostgresConfig {
	return ctx.Configuration.Spec.Services.Ledger.Postgres
}

func (l ledgerModule) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ModuleContext) modules.Services {
				return modules.Services{{
					EnvPrefix:               "NUMARY_",
					ListenEnvVar:            "SERVER_HTTP_BIND_ADDRESS",
					InjectPostgresVariables: true,
					HasVersionEndpoint:      true,
					ExposeHTTP:              true,
					NeedTopic:               true,
					Annotations:             ctx.Configuration.Spec.Services.Ledger.Annotations.Service,
					Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
						env := modules.NewEnv().Append(
							modules.Env("STORAGE_DRIVER", "postgres"),
							modules.Env("PUBLISHER_TOPIC_MAPPING", "*:"+resolveContext.Stack.GetServiceName("ledger")),
						).Append(modules.BrokerEnvVarsWithPrefix(resolveContext.Configuration.Spec.Broker, "ledger")...)

						if resolveContext.Configuration.Spec.Services.Ledger.AllowPastTimestamps {
							env = env.Append(modules.Env("COMMIT_POLICY", "allow-past-timestamps"))
						}

						// Strategy
						strategy := "memory"
						if resolveContext.Configuration.Spec.Services.Ledger.Locking.Strategy != "" {
							strategy = resolveContext.Configuration.Spec.Services.Ledger.Locking.Strategy
						}
						env = env.Append(modules.Env("LOCKING_STRATEGY", strategy))
						if strategy == "redis" {
							redisConfiguration := resolveContext.Configuration.Spec.Services.Ledger.Locking.Redis
							env = env.Append(
								modules.Env("LOCK_STRATEGY_REDIS_URL", redisConfiguration.Uri),
								modules.Env("LOCK_STRATEGY_REDIS_TLS_ENABLED", strconv.FormatBool(redisConfiguration.TLS)),
								modules.Env("LOCK_STRATEGY_REDIS_TLS_INSECURE", strconv.FormatBool(redisConfiguration.InsecureTLS)),
							)

							if redisConfiguration.Duration != 0 {
								env = append(env, modules.Env("LOCK_STRATEGY_REDIS_DURATION", redisConfiguration.Duration.String()))
							}

							if redisConfiguration.Retry != 0 {
								env = append(env, modules.Env("LOCK_STRATEGY_REDIS_RETRY", redisConfiguration.Retry.String()))
							}
						}

						return modules.Container{
							Resources: getResourcesWithDefault(
								resolveContext.Configuration.Spec.Services.Ledger.ResourceProperties,
								modules.ResourceSizeSmall(),
							),
							Image: modules.GetImage("ledger", resolveContext.Versions.Spec.Ledger),
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
							Resources: getResourcesWithDefault(
								resolveContext.Configuration.Spec.Services.Ledger.ResourceProperties,
								modules.ResourceSizeSmall(),
							),
							Image: modules.GetImage("ledger", resolveContext.Versions.Spec.Ledger),
							Env: env.Append(
								modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
							),
						}
					},
				}}
			},
		},
	}
}

var _ modules.Module = (*ledgerModule)(nil)
var _ modules.PostgresAwareModule = (*ledgerModule)(nil)

func init() {
	modules.Register("ledger", &ledgerModule{})
}
