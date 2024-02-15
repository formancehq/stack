package webhooks

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokertopicconsumers"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

func deploymentEnvVars(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumers brokertopicconsumers.Consumers) ([]v1.EnvVar, error) {

	brokerURI, err := settings.RequireURL(ctx, stack.Name, "broker.dsn")
	if err != nil {
		return nil, err
	}
	if brokerURI == nil {
		return nil, errors.New("missing broker configuration")
	}

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseKind(ctx, webhooks))
	if err != nil {
		return nil, err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	env = append(env, gatewayEnv...)
	env = append(env, core.GetDevEnvVars(stack, webhooks)...)

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "webhooks", webhooks.Spec.Auth)
	if err != nil {
		return nil, err
	}

	env = append(env, authEnvVars...)
	env = append(env, databases.GetPostgresEnvVars(database)...)
	env = append(env, settings.GetBrokerEnvVars(brokerURI, stack.Name, "webhooks")...)
	env = append(env, core.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))

	return env, nil
}

func createAPIDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumers brokertopicconsumers.Consumers, version string, withWorker bool) error {

	image, err := registries.GetImage(ctx, stack, "webhooks", version)
	if err != nil {
		return err
	}

	env, err := deploymentEnvVars(ctx, stack, webhooks, database, consumers)
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
		env = append(env, core.Env("KAFKA_TOPICS", strings.Join(Map(consumers, func(from *v1beta1.BrokerTopicConsumer) string {
			return fmt.Sprintf("%s-%s", stack.Name, from.Spec.Service)
		}), " ")))
	}

	_, err = deployments.CreateOrUpdate(ctx, stack, webhooks, "webhooks",
		deployments.WithMatchingLabels("webhooks"),
		deployments.WithServiceAccountName(database.Status.URI.Query().Get("awsRole")),
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

func createWorkerDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumers brokertopicconsumers.Consumers, version string, withWorker bool) error {

	image, err := registries.GetImage(ctx, stack, "webhooks", version)
	if err != nil {
		return err
	}

	env, err := deploymentEnvVars(ctx, stack, webhooks, database, consumers)
	if err != nil {
		return err
	}

	env = append(env, core.Env("WORKER", "true"))
	env = append(env, core.Env("KAFKA_TOPICS", strings.Join(Map(consumers, func(from *v1beta1.BrokerTopicConsumer) string {
		return fmt.Sprintf("%s-%s", stack.Name, from.Spec.Service)
	}), " ")))

	_, err = deployments.CreateOrUpdate(ctx, stack, webhooks, "webhooks-worker",
		deployments.WithMatchingLabels("webhooks-worker"),
		deployments.WithServiceAccountName(database.Status.URI.Query().Get("awsRole")),
		deployments.WithContainers(v1.Container{
			Name:  "worker",
			Env:   env,
			Image: image,
			Args:  []string{"worker"},
		}),
	)
	return err
}

func createSingleDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumers brokertopicconsumers.Consumers, version string) error {
	return createAPIDeployment(ctx, stack, webhooks, database, consumers, version, true)
}

func createDualDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumers brokertopicconsumers.Consumers, version string) error {
	if err := createAPIDeployment(ctx, stack, webhooks, database, consumers, version, false); err != nil {
		return err
	}
	if err := createWorkerDeployment(ctx, stack, webhooks, database, consumers, version, false); err != nil {
		return err
	}

	return nil
}
