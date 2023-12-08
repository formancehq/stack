package reconciliation

import (
	"github.com/formancehq/operator/apis/stack/v1beta3"
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

type module struct{}

func (o module) Name() string {
	return "reconciliation"
}

func (o module) IsEE() bool {
	return true
}

func (p module) Postgres(ctx modules.ReconciliationConfig) v1beta3.PostgresConfig {
	return ctx.Configuration.Spec.Services.Reconciliation.Postgres
}

func (o module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			DatabaseMigration: &modules.DatabaseMigration{
				Name:     "v0.0.0",
				Shutdown: false,
				Command:  []string{"migrate", "up"},
				AdditionalEnv: func(config modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return modules.Services{
					{
						InjectPostgresVariables: true,
						AuthConfiguration: func(config modules.ReconciliationConfig) stackv1beta3.ClientConfiguration {
							return stackv1beta3.NewClientConfiguration()
						},
						HasVersionEndpoint: true,
						ListenEnvVar:       "LISTEN",
						ExposeHTTP:         modules.DefaultExposeHTTP,
						Annotations:        ctx.Configuration.Spec.Services.Reconciliation.Annotations.Service,
						Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
							return modules.Container{
								Env:   reconciliationEnvVars(resolveContext),
								Image: modules.GetImage("reconciliation", resolveContext.Versions.Spec.Reconciliation),
								Resources: modules.GetResourcesWithDefault(
									resolveContext.Configuration.Spec.Services.Reconciliation.ResourceProperties,
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

func reconciliationEnvVars(resolveContext modules.ContainerResolutionConfiguration) modules.ContainerEnv {
	env := modules.NewEnv().Append(
		modules.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
		modules.Env("STACK_CLIENT_ID", resolveContext.Stack.Status.StaticAuthClients["reconciliation"].ID),
		modules.Env("STACK_CLIENT_SECRET", resolveContext.Stack.Status.StaticAuthClients["reconciliation"].Secrets[0]),
	)

	return env
}

var Module = &module{}

var _ modules.Module = Module
var _ modules.PostgresAwareModule = Module

func init() {
	modules.Register(Module)
}
