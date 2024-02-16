package orchestrations

import (
	"fmt"
	"github.com/formancehq/operator/internal/resources/resourcereferences"
	"strings"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/formancehq/operator/internal/resources/settings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
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
	consumers []*v1beta1.BrokerTopicConsumer, image string) error {

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, LowerCamelCaseKind(ctx, orchestration))
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
	env = append(env, databases.GetPostgresEnvVars(database)...)

	temporalURI, err := settings.RequireURL(ctx, stack.Name, "temporal.dsn")
	if err != nil {
		return err
	}

	if err := validateTemporalURI(temporalURI); err != nil {
		return err
	}

	var resourceReference *v1beta1.ResourceReference
	if secret := temporalURI.Query().Get("secret"); secret != "" {
		resourceReference, err = resourcereferences.Create(ctx, database, "temporal", secret, &v1.Secret{})
	} else {
		err = resourcereferences.Delete(ctx, database, "temporal")
	}
	if err != nil {
		return err
	}

	env = append(env,
		Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		Env("TEMPORAL_TASK_QUEUE", stack.Name),
		Env("TEMPORAL_ADDRESS", temporalURI.Host),
		Env("TEMPORAL_NAMESPACE", temporalURI.Path[1:]),
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

	if secret := temporalURI.Query().Get("secret"); secret == "" {
		temporalTLSCrt, err := settings.GetStringOrEmpty(ctx, stack.Name, "temporal.tls.crt")
		if err != nil {
			return err
		}

		temporalTLSKey, err := settings.GetStringOrEmpty(ctx, stack.Name, "temporal.tls.key")
		if err != nil {
			return err
		}

		env = append(env,
			Env("TEMPORAL_SSL_CLIENT_KEY", temporalTLSKey),
			Env("TEMPORAL_SSL_CLIENT_CERT", temporalTLSCrt),
		)
	} else {
		env = append(env,
			EnvFromSecret("TEMPORAL_SSL_CLIENT_KEY", secret, "tls.key"),
			EnvFromSecret("TEMPORAL_SSL_CLIENT_CERT", secret, "tls.crt"),
		)
	}

	brokerEnvVars, err := settings.ResolveBrokerEnvVars(ctx, stack, "orchestration")
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	env = append(env, brokerEnvVars...)

	_, err = deployments.CreateOrUpdate(ctx, stack, orchestration, "orchestration",
		resourcereferences.Annotate[*appsv1.Deployment]("temporal-secret-hash", resourceReference),
		deployments.WithServiceAccountName(database.Status.URI.Query().Get("awsRole")),
		deployments.WithMatchingLabels("orchestration"),
		deployments.WithContainers(v1.Container{
			Name:          "api",
			Env:           env,
			Image:         image,
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)
	return err
}

func validateTemporalURI(temporalURI *v1beta1.URI) error {
	if temporalURI.Scheme != "temporal" {
		return fmt.Errorf("invalid temporal uri: %s", temporalURI.String())
	}

	if temporalURI.Path == "" {
		return fmt.Errorf("invalid temporal uri: %s", temporalURI.String())
	}

	if !strings.HasPrefix(temporalURI.Path, "/") {
		return fmt.Errorf("invalid temporal uri: %s", temporalURI.String())
	}

	return nil
}
