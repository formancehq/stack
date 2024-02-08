package wallets

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	v1 "k8s.io/api/core/v1"
)

func createDeployment(ctx core.Context, stack *v1beta1.Stack, wallets *v1beta1.Wallets,
	authClient *v1beta1.AuthClient, version string) error {
	env := make([]v1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseKind(ctx, wallets))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	env = append(env, gatewayEnv...)
	env = append(env, core.GetDevEnvVars(stack, wallets)...)
	if authClient != nil {
		env = append(env, authclients.GetEnvVars(authClient)...)
	}

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "wallets", wallets.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	image, err := registries.GetImage(ctx, stack, "wallets", version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, stack, wallets, "wallets",
		deployments.WithContainers(v1.Container{
			Name:          "wallets",
			Args:          []string{"serve"},
			Env:           env,
			Image:         image,
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
		deployments.WithMatchingLabels("wallets"),
	)
	return err
}
