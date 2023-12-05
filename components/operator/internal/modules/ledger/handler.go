package ledger

import (
	"net/http"
	"strconv"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

type module struct{}

func (l module) Name() string {
	return "ledger"
}

func (l module) Postgres(ctx modules.ReconciliationConfig) v1beta3.PostgresConfig {
	return ctx.Configuration.Spec.Services.Ledger.Postgres
}

func (l module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return modules.Services{{
					EnvPrefix:               "NUMARY_",
					ListenEnvVar:            "SERVER_HTTP_BIND_ADDRESS",
					InjectPostgresVariables: true,
					HasVersionEndpoint:      true,
					ExposeHTTP:              modules.DefaultExposeHTTP,
					Topics:                  &modules.Topics{Name: l.Name()},
					Annotations:             ctx.Configuration.Spec.Services.Ledger.Annotations.Service,
					Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
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
						env = env.Append(modules.Env("LOCK_STRATEGY", strategy))
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
							Resources: modules.GetResourcesWithDefault(
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
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				Command:  []string{"buckets", "upgrade-all"},
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						// todo: make ledger v2 use same env vars as other services
						modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {

				createContainer := func(env modules.ContainerEnv) func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
					return func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
						env := env.Append(
							modules.Env("STORAGE_DRIVER", "postgres"),
							modules.Env("PUBLISHER_TOPIC_MAPPING", "*:"+resolveContext.Stack.GetServiceName("ledger")),
						).Append(modules.BrokerEnvVars(resolveContext.Configuration.Spec.Broker, "ledger")...)

						return modules.Container{
							Resources: modules.GetResourcesWithDefault(
								resolveContext.Configuration.Spec.Services.Ledger.ResourceProperties,
								modules.ResourceSizeSmall(),
							),
							Image: modules.GetImage("ledger", resolveContext.Versions.Spec.Ledger),
							Env: env.Append(
								modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
							),
						}
					}
				}

				main := modules.Service{
					ListenEnvVar:            "BIND",
					InjectPostgresVariables: true,
					HasVersionEndpoint:      true,
					ExposeHTTP:              modules.DefaultExposeHTTP,
					Topics:                  &modules.Topics{Name: l.Name()},
					Container:               createContainer(modules.NewEnv()),
				}
				ret := modules.Services{&main}

				if ctx.Configuration.Spec.Services.Ledger.DeploymentStrategy == v1beta3.DeploymentStrategyMonoWriterMultipleReader {
					cp := main
					cp.Name = "read"
					cp.ExposeHTTP = &modules.ExposeHTTP{
						Methods: []string{http.MethodGet, http.MethodOptions, http.MethodHead},
					}
					cp.Container = createContainer(modules.NewEnv().Append(modules.Env("READ_ONLY", "true")))
					ret = append(ret, &cp)

					main.ExposeHTTP = &modules.ExposeHTTP{
						Methods: []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
					}
					main.Replicas = pointer.For(int32(1))
				}

				return ret
			},
		},
	}
}

var Module = &module{}

var _ modules.Module = Module
var _ modules.PostgresAwareModule = Module

func init() {
	modules.Register(Module)
}
