package payments

import (
	"github.com/formancehq/operator/internal/resources/settings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/services"
	corev1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

func commonEnvVars(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) ([]v1.EnvVar, error) {
	env := make([]v1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseName(ctx, payments))
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
	env = append(env, databases.GetPostgresEnvVars(database)...)
	encryptionKey := payments.Spec.EncryptionKey
	if encryptionKey == "" {
		encryptionKey, err = settings.GetStringOrEmpty(ctx, stack.Name, "payments.encryption-key")
		if err != nil {
			return nil, err
		}
	}
	env = append(env,
		core.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
		core.Env("CONFIG_ENCRYPTION_KEY", encryptionKey),
	)

	return env, nil
}

func createFullDeployment(ctx core.Context, stack *v1beta1.Stack,
	payments *v1beta1.Payments, database *v1beta1.Database, version string) error {

	env, err := commonEnvVars(ctx, stack, payments, database)
	if err != nil {
		return err
	}

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "payments", payments.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	image, err := registries.GetImage(ctx, stack, "payments", version)
	if err != nil {
		return err
	}

	topic, err := brokertopics.Find(ctx, stack, "payments")
	if err != nil {
		return err
	}

	if topic != nil {
		if !topic.Status.Ready {
			return core.NewApplicationError("topic %s is not yet ready", topic.Name)
		}

		env = append(env, settings.GetBrokerEnvVars(topic.Status.URI, stack.Name, "payments")...)
		env = append(env, core.Env("PUBLISHER_TOPIC_MAPPING", "*:"+core.GetObjectName(stack.Name, "payments")))
	}

	_, err = deployments.CreateOrUpdate(ctx, stack, payments, "payments",
		deployments.WithMatchingLabels("payments"),
		deployments.WithContainers(v1.Container{
			Name:          "api",
			Args:          []string{"serve"},
			Env:           env,
			Image:         image,
			LivenessProbe: deployments.DefaultLiveness("http", deployments.WithProbePath("/_health")),
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
		}),
		setInitContainer(payments, database, image),
	)
	if err != nil {
		return err
	}

	return nil
}

func createReadDeployment(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database, version string) error {

	env, err := commonEnvVars(ctx, stack, payments, database)
	if err != nil {
		return err
	}

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "payments", payments.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	image, err := registries.GetImage(ctx, stack, "payments", version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, stack, payments, "payments-read",
		deployments.WithMatchingLabels("payments-read"),
		deployments.WithContainers(v1.Container{
			Name:          "api",
			Args:          []string{"api", "serve"},
			Env:           env,
			Image:         image,
			LivenessProbe: deployments.DefaultLiveness("http", deployments.WithProbePath("/_health")),
			Ports:         []v1.ContainerPort{deployments.StandardHTTPPort()},
		}),
	)
	if err != nil {
		return err
	}

	_, err = services.Create(ctx, payments, "payments-read", services.WithDefault("payments-read"))
	if err != nil {
		return err
	}

	return nil
}

func createConnectorsDeployment(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments,
	database *v1beta1.Database, version string) error {

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
			return core.NewApplicationError("topic %s is not yet ready", topic.Name)
		}

		env = append(env, settings.GetBrokerEnvVars(topic.Status.URI, stack.Name, "payments")...)
		env = append(env, core.Env("PUBLISHER_TOPIC_MAPPING", "*:"+core.GetObjectName(stack.Name, "payments")))
	}

	image, err := registries.GetImage(ctx, stack, "payments", version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, stack, payments, "payments-connectors",
		deployments.WithMatchingLabels("payments-connectors"),
		deployments.WithContainers(v1.Container{
			Name:  "connectors",
			Args:  []string{"connectors", "serve"},
			Env:   env,
			Image: image,
			Ports: []v1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http",
				deployments.WithProbePath("/_health")),
		}),
		setInitContainer(payments, database, image),
	)
	if err != nil {
		return err
	}

	_, err = services.Create(ctx, payments, "payments-connectors", services.WithDefault("payments-connectors"))
	if err != nil {
		return err
	}

	return err
}

func createGateway(ctx core.Context, stack *v1beta1.Stack, p *v1beta1.Payments) error {

	caddyfileConfigMap, err := settings.CreateCaddyfileConfigMap(ctx, stack, "payments", Caddyfile, map[string]any{
		"Debug": stack.Spec.Debug || p.Spec.Debug,
	}, core.WithController[*v1.ConfigMap](ctx.GetScheme(), p))
	if err != nil {
		return err
	}

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseName(ctx, p))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, core.GetDevEnvVars(stack, p)...)

	image, err := registries.GetImage(ctx, stack, "gateway", "latest")
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, stack, p, "payments",
		settings.ConfigureCaddy(caddyfileConfigMap, image, env),
		deployments.WithMatchingLabels("payments"),
	)
	return err
}

func setInitContainer(payments *v1beta1.Payments, database *v1beta1.Database, image string) func(t *corev1.Deployment) error {
	return func(t *corev1.Deployment) error {
		t.Spec.Template.Spec.InitContainers = []v1.Container{
			databases.MigrateDatabaseContainer(image, database,
				func(m *databases.MigrationConfiguration) {
					m.AdditionalEnv = []v1.EnvVar{
						core.Env("CONFIG_ENCRYPTION_KEY", payments.Spec.EncryptionKey),
					}
				},
			),
		}

		return nil
	}
}
