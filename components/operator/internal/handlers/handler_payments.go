package handlers

import (
	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

func init() {
	env := func(resolveContext modules.ContainerResolutionContext) modules.ContainerEnv {
		return modules.BrokerEnvVars(resolveContext.Configuration.Spec.Broker, "payments").
			Append(
				modules.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
				modules.Env("CONFIG_ENCRYPTION_KEY", resolveContext.Configuration.Spec.Services.Payments.EncryptionKey),
				modules.Env("PUBLISHER_TOPIC_MAPPING", "*:"+resolveContext.Stack.GetServiceName("payments")),
			)
	}
	modules.Register("payments", modules.Module{
		Postgres: func(ctx modules.Context) v1beta3.PostgresConfig {
			return ctx.Configuration.Spec.Services.Payments.Postgres
		},
		Versions: map[string]modules.Version{
			"v0.0.0": {
				Services: func(ctx modules.ModuleContext) modules.Services {
					migrateCommand := []string{"payments", "migrate"}
					if ctx.HasVersionLower("v0.7.0") {
						migrateCommand = append(migrateCommand, "up")
					}
					return modules.Services{{
						InjectPostgresVariables: true,
						HasVersionEndpoint:      true,
						ListenEnvVar:            "LISTEN",
						ExposeHTTP:              true,
						NeedTopic:               true,
						Liveness:                modules.LivenessLegacy,
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
							return modules.Container{
								Env:       env(resolveContext),
								Image:     modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
								Resources: modules.ResourceSizeSmall(),
							}
						},
						InitContainer: func(resolveContext modules.ContainerResolutionContext) []modules.Container {
							return []modules.Container{{
								Name:    "migrate",
								Image:   modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
								Env:     env(resolveContext),
								Command: migrateCommand,
							}}
						},
					}}
				},
			},
		},
	})
}
