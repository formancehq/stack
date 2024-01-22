package orchestrations

import (
	"fmt"
	"github.com/formancehq/operator/internal/resources/settings"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

func createAuthClient(ctx Context, stack *v1beta1.Stack, orchestration *v1beta1.Orchestration) (*v1beta1.AuthClient, error) {

	hasAuth, err := HasDependency(ctx, stack.Name, &v1beta1.Auth{})
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

func createDeployment(ctx Context, stack *v1beta1.Stack, orchestration *v1beta1.Orchestration,
	database *v1beta1.Database, client *v1beta1.AuthClient,
	consumers []*v1beta1.BrokerTopicConsumer, version string) error {

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVarsIfEnabled(ctx, stack, GetModuleName(ctx, orchestration))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	env = append(env, gatewayEnv...)
	env = append(env, GetDevEnvVars(stack, orchestration)...)
	env = append(env, databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfiguration, database.Status.Configuration.Database)...)

	temporalConfiguration, err := settings.FindTemporalConfiguration(ctx, stack)
	if err != nil {
		return err
	}

	env = append(env,
		Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		Env("TEMPORAL_TASK_QUEUE", stack.Name),
		Env("TEMPORAL_ADDRESS", temporalConfiguration.Address),
		Env("TEMPORAL_NAMESPACE", temporalConfiguration.Namespace),
		Env("WORKER", "true"),
		Env("TOPICS", strings.Join(collectionutils.Map(consumers, func(from *v1beta1.BrokerTopicConsumer) string {
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

	if temporalConfiguration.TLS.SecretName == "" {
		env = append(env,
			Env("TEMPORAL_SSL_CLIENT_KEY", temporalConfiguration.TLS.Key),
			Env("TEMPORAL_SSL_CLIENT_CERT", temporalConfiguration.TLS.CRT),
		)
	} else {
		env = append(env,
			EnvFromSecret("TEMPORAL_SSL_CLIENT_KEY", temporalConfiguration.TLS.SecretName, "tls.key"),
			EnvFromSecret("TEMPORAL_SSL_CLIENT_CERT", temporalConfiguration.TLS.SecretName, "tls.crt"),
		)
	}

	brokerEnvVars, err := settings.ResolveBrokerEnvVars(ctx, stack, "orchestration")
	if err != nil && !errors.Is(err, ErrNotFound) {
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
			Resources:     GetResourcesRequirementsWithDefault(orchestration.Spec.ResourceRequirements, ResourceSizeSmall()),
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)
	return err
}
