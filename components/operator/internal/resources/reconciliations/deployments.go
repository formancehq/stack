package reconciliations

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/internal/resources/registries"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

func createDeployment(ctx core.Context, stack *v1beta1.Stack, reconciliation *v1beta1.Reconciliation,
	database *v1beta1.Database, authClient *v1beta1.AuthClient, version string) error {
	env := make([]v1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, core.GetModuleName(ctx, reconciliation))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	env = append(env, gatewayEnv...)
	env = append(env, core.GetDevEnvVars(stack, reconciliation)...)

	env = append(env,
		databases.PostgresEnvVars(
			database.Status.Configuration.DatabaseConfigurationSpec,
			core.GetObjectName(stack.Name, "reconciliation"),
		)...,
	)
	env = append(env, core.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"))
	env = append(env, authclients.GetEnvVars(authClient)...)

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "reconciliation", reconciliation.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	image, err := registries.GetImage(ctx, stack, "reconciliation", version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, reconciliation, "reconciliation",
		func(t *v12.Deployment) {
			t.Spec.Template.Spec.Containers = []v1.Container{{
				Name:          "reconciliation",
				Env:           env,
				Image:         image,
				Resources:     core.GetResourcesRequirementsWithDefault(reconciliation.Spec.ResourceRequirements, core.ResourceSizeSmall()),
				Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
				LivenessProbe: deployments.DefaultLiveness("http"),
			}}
		},
		deployments.WithMatchingLabels("reconciliation"),
	)
	return err
}
