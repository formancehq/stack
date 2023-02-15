package handlers

import (
	"strings"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

func init() {
	modules.Register("webhooks", modules.Module{
		Postgres: func(ctx modules.Context) stackv1beta3.PostgresConfig {
			return ctx.Configuration.Spec.Services.Webhooks.Postgres
		},
		Services: func(ctx modules.Context) modules.Services {
			return modules.Services{
				modules.Service{
					HasVersionEndpoint:      true,
					Port:                    8080,
					InjectPostgresVariables: true,
					Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
						return modules.Container{
							Image: modules.GetImage("webhooks", resolveContext.Versions.Spec.Webhooks),
							Env:   webhooksEnvVars(resolveContext.Configuration),
						}
					},
				},
				modules.Service{
					Name:                    "worker",
					InjectPostgresVariables: true,
					Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
						return modules.Container{
							Image: modules.GetImage("webhooks", resolveContext.Versions.Spec.Webhooks),
							Env: webhooksEnvVars(resolveContext.Configuration).Append(
								modules.Env("KAFKA_TOPICS", strings.Join([]string{
									resolveContext.Stack.GetServiceName("ledger"),
									resolveContext.Stack.GetServiceName("payments"),
								}, " ")),
							),
							Command:  []string{"/usr/bin/webhooks", "worker"},
							Liveness: modules.LivenessDisable,
						}
					},
				},
			}
		},
	})
}

func webhooksEnvVars(configuration *stackv1beta3.Configuration) modules.ContainerEnv {
	return modules.BrokerEnvVars(configuration.Spec.Broker, "webhooks").
		Append(modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))
}
