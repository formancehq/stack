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
		Versions: map[string]modules.Version{
			"v0.0.0": {
				Services: func(ctx modules.ModuleContext) modules.Services {
					return modules.Services{
						{
							HasVersionEndpoint:      true,
							ExposeHTTP:              true,
							InjectPostgresVariables: true,
							ListenEnvVar:            "LISTEN",
							Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
								return modules.Container{
									Image:     modules.GetImage("webhooks", resolveContext.Versions.Spec.Webhooks),
									Env:       webhooksEnvVars(resolveContext.Configuration),
									Resources: modules.ResourceSizeSmall(),
								}
							},
						},
						{
							Name:                    "worker",
							InjectPostgresVariables: true,
							ListenEnvVar:            "LISTEN",
							Liveness:                modules.LivenessDisable,
							Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
								return modules.Container{
									Image: modules.GetImage("webhooks", resolveContext.Versions.Spec.Webhooks),
									Env: webhooksEnvVars(resolveContext.Configuration).Append(
										modules.Env("KAFKA_TOPICS", strings.Join([]string{
											resolveContext.Stack.GetServiceName("ledger"),
											resolveContext.Stack.GetServiceName("payments"),
										}, " ")),
									),
									Args:      []string{"worker"},
									Resources: modules.ResourceSizeSmall(),
								}
							},
						},
					}
				},
			},
		},
	})
}

func webhooksEnvVars(configuration *stackv1beta3.Configuration) modules.ContainerEnv {
	return modules.BrokerEnvVars(configuration.Spec.Broker, "webhooks").
		Append(modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))
}
