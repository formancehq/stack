package orchestration

import (
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

type module struct{}

func (o module) Postgres(ctx modules.Context) stackv1beta3.PostgresConfig {
	return ctx.Configuration.Spec.Services.Orchestration.Postgres
}

func (o module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ModuleContext) modules.Services {
				return modules.Services{
					{
						Port: 8080,
						AuthConfiguration: func(resolveContext modules.ModuleContext) stackv1beta3.ClientConfiguration {
							return stackv1beta3.NewClientConfiguration()
						},
						ExposeHTTP:              true,
						HasVersionEndpoint:      true,
						InjectPostgresVariables: true,
						Annotations:             ctx.Configuration.Spec.Services.Orchestration.Annotations.Service,
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
							return modules.Container{
								Env:   orchestrationEnvVars(resolveContext),
								Image: modules.GetImage("orchestration", resolveContext.Versions.Spec.Orchestration),
								Resources: modules.GetResourcesWithDefault(
									resolveContext.Configuration.Spec.Services.Orchestration.ResourceProperties,
									modules.ResourceSizeSmall(),
								),
							}
						},
					},
					{
						Name:                    "worker",
						InjectPostgresVariables: true,
						Liveness:                modules.LivenessDisable,
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
							return modules.Container{
								Env:   orchestrationEnvVars(resolveContext),
								Image: modules.GetImage("orchestration", resolveContext.Versions.Spec.Orchestration),
								Args:  []string{"worker"},
								Resources: modules.GetResourcesWithDefault(
									resolveContext.Configuration.Spec.Services.Orchestration.ResourceProperties,
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

var _ modules.Module = (*module)(nil)
var _ modules.PostgresAwareModule = (*module)(nil)

func init() {
	modules.Register("orchestration", &module{})
}

func orchestrationEnvVars(resolveContext modules.ContainerResolutionContext) modules.ContainerEnv {
	return modules.NewEnv().Append(
		modules.Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		modules.Env("TEMPORAL_TASK_QUEUE", resolveContext.Stack.Name),
		modules.Env("TEMPORAL_ADDRESS", resolveContext.Configuration.Spec.Temporal.Address),
		modules.Env("TEMPORAL_NAMESPACE", resolveContext.Configuration.Spec.Temporal.Namespace),
		modules.Env("TEMPORAL_SSL_CLIENT_KEY", resolveContext.Configuration.Spec.Temporal.TLS.Key),
		modules.Env("TEMPORAL_SSL_CLIENT_CERT", resolveContext.Configuration.Spec.Temporal.TLS.CRT),
		modules.Env("STACK_CLIENT_ID", resolveContext.Stack.Status.StaticAuthClients["orchestration"].ID),
		modules.Env("STACK_CLIENT_SECRET", resolveContext.Stack.Status.StaticAuthClients["orchestration"].Secrets[0]),
	)
}
