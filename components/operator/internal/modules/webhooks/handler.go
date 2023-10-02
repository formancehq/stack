package webhooks

import (
	"strings"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

type module struct{}

func (w module) Name() string {
	return "webhooks"
}

func (w module) Postgres(ctx modules.ReconciliationConfig) stackv1beta3.PostgresConfig {
	return ctx.Configuration.Spec.Services.Webhooks.Postgres
}

func (w module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return modules.Services{
					{
						HasVersionEndpoint:      true,
						ExposeHTTP:              modules.DefaultExposeHTTP,
						InjectPostgresVariables: true,
						ListenEnvVar:            "LISTEN",
						Annotations:             ctx.Configuration.Spec.Services.Webhooks.Annotations.Service,
						Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
							return modules.Container{
								Image: modules.GetImage("webhooks", resolveContext.Versions.Spec.Webhooks),
								Env:   webhooksEnvVars(resolveContext.Configuration),
								Resources: modules.GetResourcesWithDefault(
									resolveContext.Configuration.Spec.Services.Webhooks.ResourceProperties,
									modules.ResourceSizeSmall(),
								),
							}
						},
					},
					{
						Name:                    "worker",
						InjectPostgresVariables: true,
						ListenEnvVar:            "LISTEN",
						Liveness:                modules.LivenessDisable,
						Annotations:             ctx.Configuration.Spec.Services.Webhooks.Annotations.Service,
						Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
							return modules.Container{
								Image: modules.GetImage("webhooks", resolveContext.Versions.Spec.Webhooks),
								Env: webhooksEnvVars(resolveContext.Configuration).Append(
									modules.Env("KAFKA_TOPICS", strings.Join([]string{
										resolveContext.Stack.GetServiceName("ledger"),
										resolveContext.Stack.GetServiceName("payments"),
									}, " ")),
								),
								Args: []string{"worker"},
								Resources: modules.GetResourcesWithDefault(
									resolveContext.Configuration.Spec.Services.Webhooks.ResourceProperties,
									modules.ResourceSizeSmall(),
								),
							}
						},
					},
				}
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

func webhooksEnvVars(configuration *stackv1beta3.Configuration) modules.ContainerEnv {
	return modules.BrokerEnvVars(configuration.Spec.Broker, "webhooks").
		Append(modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))
}
