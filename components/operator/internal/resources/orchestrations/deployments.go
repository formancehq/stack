package orchestrations

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

func createAuthClient(ctx core.Context, stack *v1beta1.Stack, orchestration *v1beta1.Orchestration) (*v1beta1.AuthClient, error) {

	hasAuth, err := core.HasDependency(ctx, stack.Name, &v1beta1.Auth{})
	if err != nil {
		return nil, err
	}
	if !hasAuth {
		return nil, nil
	}

	return authclients.Create(ctx, stack, orchestration, "orchestration",
		func(spec *v1beta1.AuthClientSpec) {
			spec.Scopes = []string{
				"ledger:read",
				"ledger:write",
				"payments:read",
				"payments:write",
				"wallets:read",
				"wallets:write",
			}
		})
}

func createDeployment(ctx core.Context, stack *v1beta1.Stack, orchestration *v1beta1.Orchestration,
	database *v1beta1.Database, client *v1beta1.AuthClient,
	consumers []*v1beta1.BrokerTopicConsumer, version string) error {

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, core.GetModuleName(ctx, orchestration))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	env = append(env, gatewayEnv...)
	env = append(env, core.GetDevEnvVars(stack, orchestration)...)
	env = append(env, databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)...)

	temporalConfiguration, err := core.RequireConfigurationObject[*v1beta1.TemporalConfiguration](ctx, stack.Name)
	if err != nil {
		return err
	}

	env = append(env,
		core.Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		core.Env("TEMPORAL_TASK_QUEUE", stack.Name),
		core.Env("TEMPORAL_ADDRESS", temporalConfiguration.Spec.Address),
		core.Env("TEMPORAL_NAMESPACE", temporalConfiguration.Spec.Namespace),
		core.Env("WORKER", "true"),
		core.Env("TOPICS", strings.Join(collectionutils.Map(consumers, func(from *v1beta1.BrokerTopicConsumer) string {
			return fmt.Sprintf("%s-%s", stack.Name, from.Spec.Service)
		}), " ")),
	)

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "orchestration", orchestration.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	if client != nil {
		env = append(env, authclients.GetEnvVars(client)...)
	}

	if temporalConfiguration.Spec.TLS.SecretName == "" {
		env = append(env,
			core.Env("TEMPORAL_SSL_CLIENT_KEY", temporalConfiguration.Spec.TLS.Key),
			core.Env("TEMPORAL_SSL_CLIENT_CERT", temporalConfiguration.Spec.TLS.CRT),
		)
	} else {
		env = append(env,
			core.EnvFromSecret("TEMPORAL_SSL_CLIENT_KEY", temporalConfiguration.Spec.TLS.SecretName, "tls.key"),
			core.EnvFromSecret("TEMPORAL_SSL_CLIENT_CERT", temporalConfiguration.Spec.TLS.SecretName, "tls.crt"),
		)
	}

	brokerEnvVars, err := brokerconfigurations.GetEnvVars(ctx, stack.Name, "orchestration")
	if err != nil && !errors.Is(err, core.ErrNotFound) {
		return err
	}
	env = append(env, brokerEnvVars...)

	image, err := registries.GetImage(ctx, stack, "orchestration", version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, orchestration, "orchestration",
		deployments.WithMatchingLabels("orchestration"),
		deployments.WithContainers(v1.Container{
			Name:          "api",
			Env:           env,
			Image:         image,
			Resources:     core.GetResourcesRequirementsWithDefault(orchestration.Spec.ResourceRequirements, core.ResourceSizeSmall()),
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)
	return err
}
