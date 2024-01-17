package webhooks

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/internal/resources/brokertopicconsumers"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	v1 "k8s.io/api/core/v1"
)

func createDeployment(ctx core.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks, database *v1beta1.Database, consumers brokertopicconsumers.Consumers) error {

	brokerConfiguration, err := core.RequireLabelledConfig[*v1beta1.BrokerConfiguration](ctx, stack.Name)
	if err != nil {
		return err
	}

	image, err := registries.GetImage(ctx, stack, "webhooks", webhooks.Spec.Version)
	if err != nil {
		return err
	}

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, core.GetModuleName(webhooks))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	env = append(env, gatewayEnv...)
	env = append(env, core.GetDevEnvVars(stack, webhooks)...)

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "webhooks", webhooks.Spec.Auth)
	if err != nil {
		return err
	}

	env = append(env, authEnvVars...)
	env = append(env, databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)...)
	env = append(env, brokerconfigurations.BrokerEnvVars(brokerConfiguration.Spec, stack.Name, "webhooks")...)
	env = append(env, core.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))
	env = append(env, core.Env("WORKER", "true"))
	env = append(env, core.Env("KAFKA_TOPICS", strings.Join(collectionutils.Map(consumers, func(from *v1beta1.BrokerTopicConsumer) string {
		return fmt.Sprintf("%s-%s", stack.Name, from.Spec.Service)
	}), " ")))

	_, err = deployments.CreateOrUpdate(ctx, webhooks, "webhooks",
		deployments.WithMatchingLabels("webhooks"),
		deployments.WithContainers(v1.Container{
			Name:          "api",
			Env:           env,
			Image:         image,
			Args:          []string{"serve", "--auto-migrate"},
			Resources:     core.GetResourcesRequirementsWithDefault(webhooks.Spec.ResourceRequirements, core.ResourceSizeSmall()),
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)
	return err
}
