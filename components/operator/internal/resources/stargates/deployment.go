package stargates

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/applications"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createDeployment(ctx core.Context, stack *v1beta1.Stack, stargate *v1beta1.Stargate, version string) error {

	env := make([]v1.EnvVar, 0)

	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseKind(ctx, stargate))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	env = append(env, gatewayEnv...)

	env = append(env, core.GetDevEnvVars(stack, stargate)...)
	env = append(env,
		core.Env("ORGANIZATION_ID", stargate.Spec.OrganizationID),
		core.Env("STACK_ID", stargate.Spec.StackID),
		core.Env("STARGATE_SERVER_URL", stargate.Spec.ServerURL),
		core.Env("GATEWAY_URL", "http://gateway:8080"),
		core.Env("STARGATE_AUTH_CLIENT_ID", stargate.Spec.Auth.ClientID),
		core.Env("STARGATE_AUTH_CLIENT_SECRET", stargate.Spec.Auth.ClientSecret),
		core.Env("STARGATE_AUTH_ISSUER_URL", stargate.Spec.Auth.Issuer),
	)

	image, err := registries.GetImage(ctx, stack, "stargate", version)
	if err != nil {
		return err
	}

	return applications.
		New(stargate, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "stargate",
			},
			Spec: appsv1.DeploymentSpec{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						Containers: []v1.Container{{
							Name:          "stargate",
							Env:           env,
							Image:         image,
							Ports:         []v1.ContainerPort{applications.StandardHTTPPort()},
							LivenessProbe: applications.DefaultLiveness("http"),
						}},
					},
				},
			},
		}).
		IsEE().
		Install(ctx)
}
