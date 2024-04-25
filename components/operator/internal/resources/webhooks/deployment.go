package webhooks

import (
	"fmt"
	"github.com/formancehq/operator/internal/resources/brokers"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/licence"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/resourcereferences"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

func deploymentEnvVars(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database) (*v1beta1.ResourceReference, []v1.EnvVar, error) {

	brokerURI, err := settings.RequireURL(ctx, stack.Name, "broker", "dsn")
	if err != nil {
		return nil, nil, err
	}
	if brokerURI == nil {
		return nil, nil, errors.New("missing broker configuration")
	}

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseKind(ctx, webhooks))
	if err != nil {
		return nil, nil, err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return nil, nil, err
	}
	env = append(env, gatewayEnv...)
	resourceReference, licenceEnvVars, err := licence.GetLicenceEnvVars(ctx, stack, "webhooks", webhooks)
	if err != nil {
		return nil, nil, err
	}
	env = append(env, licenceEnvVars...)
	env = append(env, core.GetDevEnvVars(stack, webhooks)...)

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "webhooks", webhooks.Spec.Auth)
	if err != nil {
		return nil, nil, err
	}

	postgresEnvVar, err := databases.GetPostgresEnvVars(ctx, stack, database)
	if err != nil {
		return nil, nil, err
	}

	brokerEnvVar, err := brokers.GetBrokerEnvVars(ctx, brokerURI, stack.Name, "webhooks")
	if err != nil {
		return nil, nil, err
	}

	env = append(env, authEnvVars...)
	env = append(env, postgresEnvVar...)
	env = append(env, brokerEnvVar...)
	env = append(env, core.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))

	return resourceReference, env, nil
}

func createAPIDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumer *v1beta1.BrokerConsumer, version string, withWorker bool) error {

	image, err := registries.GetImage(ctx, stack, "webhooks", version)
	if err != nil {
		return err
	}

	resourceReference, env, err := deploymentEnvVars(ctx, stack, webhooks, database)
	if err != nil {
		return err
	}

	args := []string{"serve"}

	// notes(gfyrag): upgrade command introduced in version v2.0.0-rc.5
	if core.IsGreaterOrEqual(version, "v2.0.0-alpha") && core.IsLower(version, "v2.0.0-rc.5") {
		args = append(args, "--auto-migrate")
	}
	if withWorker {
		env = append(env, core.Env("WORKER", "true"))

		topics, err := brokers.GetTopicsEnvVars(ctx, stack, "KAFKA_TOPICS", consumer.Spec.Services...)
		if err != nil {
			return err
		}
		env = append(env, topics...)
	}

	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, webhooks, "webhooks",
		resourcereferences.Annotate("licence-secret-hash", resourceReference),
		deployments.WithReplicasFromSettings(ctx, stack),
		deployments.WithMatchingLabels("webhooks"),
		deployments.WithServiceAccountName(serviceAccountName),
		deployments.WithContainers(v1.Container{
			Name:          "api",
			Env:           env,
			Image:         image,
			Args:          args,
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)
	return err
}

func createWorkerDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumer *v1beta1.BrokerConsumer, version string) error {

	image, err := registries.GetImage(ctx, stack, "webhooks", version)
	if err != nil {
		return err
	}

	resourceReference, env, err := deploymentEnvVars(ctx, stack, webhooks, database)
	if err != nil {
		return err
	}

	env = append(env, core.Env("WORKER", "true"))
	env = append(env, core.Env("KAFKA_TOPICS", strings.Join(Map(consumer.Spec.Services, func(from string) string {
		return fmt.Sprintf("%s-%s", stack.Name, from)
	}), " ")))

	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, webhooks, "webhooks-worker",
		resourcereferences.Annotate("licence-secret-hash", resourceReference),
		deployments.WithMatchingLabels("webhooks-worker"),
		deployments.WithServiceAccountName(serviceAccountName),
		deployments.WithContainers(v1.Container{
			Name:  "worker",
			Env:   env,
			Image: image,
			Args:  []string{"worker"},
		}),
	)
	return err
}

func createSingleDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumer *v1beta1.BrokerConsumer, version string) error {
	return createAPIDeployment(ctx, stack, webhooks, database, consumer, version, true)
}

func createDualDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumer *v1beta1.BrokerConsumer, version string) error {
	if err := createAPIDeployment(ctx, stack, webhooks, database, consumer, version, false); err != nil {
		return err
	}
	if err := createWorkerDeployment(ctx, stack, webhooks, database, consumer, version); err != nil {
		return err
	}

	return nil
}
