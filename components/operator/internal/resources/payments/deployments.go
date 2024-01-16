package payments

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/services"
	v12 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
)

func commonEnvVars(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) ([]v1.EnvVar, error) {
	env := make([]v1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, core.GetModuleName(payments))
	if err != nil {
		return nil, err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	env = append(env, gatewayEnv...)
	env = append(env, core.GetDevEnvVars(stack, payments)...)
	env = append(env, databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)...)
	env = append(env,
		core.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
		core.Env("CONFIG_ENCRYPTION_KEY", payments.Spec.EncryptionKey),
	)

	return env, nil
}

func createFullDeployment(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) error {

	env, err := commonEnvVars(ctx, stack, payments, database)
	if err != nil {
		return err
	}

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "payments", payments.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	image, err := registries.GetImage(ctx, stack, "payments", payments.Spec.Version)
	if err != nil {
		return err
	}

	topic, err := brokertopics.Find(ctx, stack, "payments")
	if err != nil {
		return err
	}

	if topic != nil {
		if !topic.Status.Ready {
			return fmt.Errorf("topic %s is not yet ready", topic.Name)
		}

		env = append(env, brokerconfigurations.BrokerEnvVars(*topic.Status.Configuration, stack.Name, "payments")...)
		env = append(env, core.Env("PUBLISHER_TOPIC_MAPPING", "*:"+core.GetObjectName(stack.Name, "payments")))
	}

	_, err = deployments.CreateOrUpdate(ctx, payments, "payments",
		deployments.WithMatchingLabels("payments"),
		deployments.WithContainers(v1.Container{
			Name:          "api",
			Args:          []string{"serve"},
			Env:           env,
			Image:         image,
			Resources:     core.GetResourcesRequirementsWithDefault(payments.Spec.ResourceRequirements, core.ResourceSizeSmall()),
			LivenessProbe: deployments.DefaultLiveness("http", deployments.WithProbePath("/_health")),
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
		}),
	)
	if err != nil {
		return err
	}

	return nil
}

func createReadDeployment(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) error {

	env, err := commonEnvVars(ctx, stack, payments, database)
	if err != nil {
		return err
	}

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "payments", payments.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	image, err := registries.GetImage(ctx, stack, "payments", payments.Spec.Version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, payments, "payments-read",
		deployments.WithMatchingLabels("payments-read"),
		deployments.WithContainers(v1.Container{
			Name:          "api",
			Args:          []string{"api", "serve"},
			Env:           env,
			Image:         image,
			Resources:     core.GetResourcesRequirementsWithDefault(payments.Spec.ResourceRequirements, core.ResourceSizeSmall()),
			LivenessProbe: deployments.DefaultLiveness("http", deployments.WithProbePath("/_health")),
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
		}),
	)
	if err != nil {
		return err
	}

	_, err = services.Create(ctx, payments, "payments-read")
	if err != nil {
		return err
	}

	return nil
}

func createConnectorsDeployment(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) error {

	env, err := commonEnvVars(ctx, stack, payments, database)
	if err != nil {
		return err
	}

	topic, err := brokertopics.Find(ctx, stack, "payments")
	if err != nil {
		return err
	}

	if topic != nil {
		if !topic.Status.Ready {
			return fmt.Errorf("topic %s is not yet ready", topic.Name)
		}

		env = append(env, brokerconfigurations.BrokerEnvVars(*topic.Status.Configuration, stack.Name, "payments")...)
		env = append(env, core.Env("PUBLISHER_TOPIC_MAPPING", "*:"+core.GetObjectName(stack.Name, "payments")))
	}

	image, err := registries.GetImage(ctx, stack, "payments", payments.Spec.Version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, payments, "payments-connectors",
		deployments.WithMatchingLabels("payments-connectors"),
		deployments.WithContainers(v1.Container{
			Name:      "connectors",
			Args:      []string{"connectors", "serve"},
			Env:       env,
			Image:     image,
			Resources: core.GetResourcesRequirementsWithDefault(payments.Spec.ResourceRequirements, core.ResourceSizeSmall()),
			Ports:     []v1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http",
				deployments.WithProbePath("/_health")),
		}),
		setInitContainer(payments, database, image),
	)
	if err != nil {
		return err
	}

	_, err = services.Create(ctx, payments, "payments-connectors")
	if err != nil {
		return err
	}

	return err
}

func createGateway(ctx core.Context, stack *v1beta1.Stack, p *v1beta1.Payments) error {

	caddyfileConfigMap, err := core.CreateCaddyfileConfigMap(ctx, stack, "payments", Caddyfile, map[string]any{
		"Debug": stack.Spec.Debug || p.Spec.Debug,
	}, core.WithController[*v1.ConfigMap](ctx.GetScheme(), p))
	if err != nil {
		return err
	}

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, core.GetModuleName(p))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, core.GetDevEnvVars(stack, p)...)

	_, err = deployments.CreateOrUpdate(ctx, p, "payments",
		core.ConfigureCaddy(caddyfileConfigMap, "caddy:2.7.6-alpine", env, nil),
		deployments.WithMatchingLabels("payments"),
	)
	return err
}

func setInitContainer(payments *v1beta1.Payments, database *v1beta1.Database, image string) func(t *v12.Deployment) {
	return func(t *v12.Deployment) {
		t.Spec.Template.Spec.InitContainers = []v1.Container{
			databases.MigrateDatabaseContainer(
				image,
				database.Status.Configuration.DatabaseConfigurationSpec,
				database.Status.Configuration.Database,
				func(m *databases.MigrationConfiguration) {
					m.AdditionalEnv = []v1.EnvVar{
						core.Env("CONFIG_ENCRYPTION_KEY", payments.Spec.EncryptionKey),
					}
				},
			),
		}
	}
}
