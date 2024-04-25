package orchestrations

import (
	"fmt"
	"github.com/formancehq/operator/internal/resources/brokers"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/licence"
	"github.com/formancehq/operator/internal/resources/resourcereferences"
	"github.com/formancehq/operator/internal/resources/settings"
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
	consumer *v1beta1.BrokerConsumer, image string) error {

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

	postgresEnvVar, err := databases.GetPostgresEnvVars(ctx, stack, database)
	if err != nil {
		return err
	}

	licenceResourceReference, licenceEnvVars, err := licence.GetLicenceEnvVars(ctx, stack, "orchestration", orchestration)
	if err != nil {
		return err
	}

	env = append(env, gatewayEnv...)
	env = append(env, licenceEnvVars...)
	env = append(env, GetDevEnvVars(stack, orchestration)...)
	env = append(env, postgresEnvVar...)

	temporalURI, err := settings.RequireURL(ctx, stack.Name, "temporal", "dsn")
	if err != nil {
		return err
	}

	if err := validateTemporalURI(temporalURI); err != nil {
		return err
	}

	var databaseResourceReference *v1beta1.ResourceReference
	if secret := temporalURI.Query().Get("secret"); secret != "" {
		databaseResourceReference, err = resourcereferences.Create(ctx, database, "temporal", secret, &v1.Secret{})
	} else {
		err = resourcereferences.Delete(ctx, database, "temporal")
	}
	if err != nil {
		return err
	}

	topics, err := brokers.GetTopicsEnvVars(ctx, stack, "TOPICS", consumer.Spec.Services...)
	if err != nil {
		return err
	}
	env = append(env, topics...)

	env = append(env,
		Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		Env("TEMPORAL_TASK_QUEUE", stack.Name),
		Env("TEMPORAL_ADDRESS", temporalURI.Host),
		Env("TEMPORAL_NAMESPACE", temporalURI.Path[1:]),
		Env("WORKER", "true"),
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
		temporalTLSCrt, err := settings.GetStringOrEmpty(ctx, stack.Name, "temporal", "tls", "crt")
		if err != nil {
			return err
		}

		temporalTLSKey, err := settings.GetStringOrEmpty(ctx, stack.Name, "temporal", "tls", "key")
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

	brokerEnvVars, err := brokers.ResolveBrokerEnvVars(ctx, stack, "orchestration")
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	env = append(env, brokerEnvVars...)

	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, orchestration, "orchestration",
		resourcereferences.Annotate("temporal-secret-hash", databaseResourceReference),
		resourcereferences.Annotate("licence-secret-hash", licenceResourceReference),
		deployments.WithServiceAccountName(serviceAccountName),
		deployments.WithReplicasFromSettings(ctx, stack),
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
