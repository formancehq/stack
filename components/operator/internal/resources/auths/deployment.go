package auths

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/internal/resources/registries"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func createDeployment(ctx Context, stack *v1beta1.Stack, auth *v1beta1.Auth, database *v1beta1.Database,
	configMap *corev1.ConfigMap, version string) (*appsv1.Deployment, error) {

	env := make([]corev1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, GetModuleName(ctx, auth))
	if err != nil {
		return nil, err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	env = append(env, gatewayEnv...)
	env = append(env, GetDevEnvVars(stack, auth)...)

	env = append(env,
		databases.PostgresEnvVars(
			database.Status.Configuration.DatabaseConfigurationSpec,
			GetObjectName(stack.Name, "auth"),
		)...,
	)
	env = append(env, Env("CONFIG", "/config/config.yaml"))

	authUrl, err := url(ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	env = append(env, Env("BASE_URL", authUrl))

	if auth.Spec.SigningKey != "" && auth.Spec.SigningKeyFromSecret != nil {
		return nil, fmt.Errorf("cannot specify signing key using both .spec.signingKey and .spec.signingKeyFromSecret fields")
	}
	if auth.Spec.SigningKey != "" {
		env = append(env, Env("SIGNING_KEY", auth.Spec.SigningKey))
	}
	if auth.Spec.SigningKeyFromSecret != nil {
		env = append(env, corev1.EnvVar{
			Name: "SIGNING_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: auth.Spec.SigningKeyFromSecret,
			},
		})
	}
	if auth.Spec.DelegatedOIDCServer != nil {
		env = append(env,
			Env("DELEGATED_CLIENT_SECRET", auth.Spec.DelegatedOIDCServer.ClientSecret),
			Env("DELEGATED_CLIENT_ID", auth.Spec.DelegatedOIDCServer.ClientID),
			Env("DELEGATED_ISSUER", auth.Spec.DelegatedOIDCServer.Issuer),
		)
	}
	if stack.Spec.Dev || auth.Spec.Dev {
		env = append(env, Env("CAOS_OIDC_DEV", "1"))
	}

	image, err := registries.GetImage(ctx, stack, "auth", version)
	if err != nil {
		return nil, err
	}

	return deployments.CreateOrUpdate(ctx, auth, "auth",
		deployments.WithMatchingLabels("auth"),
		func(t *appsv1.Deployment) {
			t.Spec.Template.Annotations = MergeMaps(t.Spec.Template.Annotations, map[string]string{
				"config-hash": HashFromConfigMaps(configMap),
			})
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:      "auth",
				Args:      []string{"serve"},
				Env:       env,
				Image:     image,
				Resources: GetResourcesRequirementsWithDefault(auth.Spec.ResourceRequirements, ResourceSizeSmall()),
				VolumeMounts: []corev1.VolumeMount{
					NewVolumeMount("config", "/config"),
				},
				Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
				LivenessProbe: deployments.DefaultLiveness("http"),
			}}
			t.Spec.Template.Spec.Volumes = []corev1.Volume{
				NewVolumeFromConfigMap("config", configMap),
			}
		},
	)
}

func url(ctx Context, stackName string) (string, error) {
	gateway := &v1beta1.Gateway{}
	ok, err := GetIfEnabled(ctx, stackName, gateway)
	if err != nil {
		return "", err
	}

	if ok {
		return gateways.URL(gateway) + "/api/auth", nil
	}

	return "http://auth:8080", nil
}