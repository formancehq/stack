package payments

import (
	"github.com/formancehq/operator/internal/resources/brokers"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/caddy"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/applications"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/services"
	v1 "k8s.io/api/core/v1"
)

func getEncryptionKey(ctx core.Context, payments *v1beta1.Payments) (string, error) {
	encryptionKey := payments.Spec.EncryptionKey
	if encryptionKey == "" {
		return settings.GetStringOrEmpty(ctx, payments.Spec.Stack, "payments", "encryption-key")
	}
	return "", nil
}

func commonEnvVars(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) ([]v1.EnvVar, error) {
	env := make([]v1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseKind(ctx, payments))
	if err != nil {
		return nil, err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return nil, err
	}

	postgresEnvVar, err := databases.GetPostgresEnvVars(ctx, stack, database)
	if err != nil {
		return nil, err
	}

	env = append(env, gatewayEnv...)
	env = append(env, core.GetDevEnvVars(stack, payments)...)
	env = append(env, postgresEnvVar...)

	encryptionKey, err := getEncryptionKey(ctx, payments)
	if err != nil {
		return nil, err
	}
	env = append(env,
		core.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
		core.Env("CONFIG_ENCRYPTION_KEY", encryptionKey),
	)

	return env, nil
}

func createFullDeployment(ctx core.Context, stack *v1beta1.Stack,
	payments *v1beta1.Payments, database *v1beta1.Database, image string) error {

	env, err := commonEnvVars(ctx, stack, payments, database)
	if err != nil {
		return err
	}

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "payments", payments.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	var broker *v1beta1.Broker
	if t, err := brokertopics.Find(ctx, stack, "payments"); err != nil {
		return err
	} else if t != nil && t.Status.Ready {
		broker = &v1beta1.Broker{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: stack.Name,
		}, broker); err != nil {
			return err
		}
	}

	if broker != nil {
		if !broker.Status.Ready {
			return core.NewPendingError().WithMessage("broker not ready")
		}
		brokerEnvVar, err := brokers.GetBrokerEnvVars(ctx, broker.Status.URI, stack.Name, "payments")
		if err != nil {
			return err
		}

		env = append(env, brokerEnvVar...)
		env = append(env, brokers.GetPublisherEnvVars(stack, broker, "payments", "")...)
	}

	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	err = applications.
		New(payments, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "payments",
			},
			Spec: appsv1.DeploymentSpec{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						ServiceAccountName: serviceAccountName,
						Containers: []v1.Container{{
							Name:          "api",
							Args:          []string{"serve"},
							Env:           env,
							Image:         image,
							LivenessProbe: applications.DefaultLiveness("http", applications.WithProbePath("/_health")),
							Ports:         []v1.ContainerPort{applications.StandardHTTPPort()},
						}},
						// Ensure empty
						InitContainers: []v1.Container{},
					},
				},
			},
		}).
		Install(ctx)
	if err != nil {
		return err
	}

	return nil
}

func createReadDeployment(ctx core.Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database, image string) error {

	env, err := commonEnvVars(ctx, stack, payments, database)
	if err != nil {
		return err
	}

	authEnvVars, err := auths.ProtectedEnvVars(ctx, stack, "payments", payments.Spec.Auth)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	err = applications.
		New(payments, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "payments-read",
			},
			Spec: appsv1.DeploymentSpec{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						ServiceAccountName: serviceAccountName,
						Containers: []v1.Container{{
							Name:          "api",
							Args:          []string{"api", "serve"},
							Env:           env,
							Image:         image,
							LivenessProbe: applications.DefaultLiveness("http", applications.WithProbePath("/_health")),
							Ports:         []v1.ContainerPort{applications.StandardHTTPPort()},
						}},
						// Ensure empty
						InitContainers: []v1.Container{},
					},
				},
			},
		}).
		Install(ctx)
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
	database *v1beta1.Database, image string) error {

	env, err := commonEnvVars(ctx, stack, payments, database)
	if err != nil {
		return err
	}

	var broker *v1beta1.Broker
	if t, err := brokertopics.Find(ctx, stack, "payments"); err != nil {
		return err
	} else if t != nil && t.Status.Ready {
		broker = &v1beta1.Broker{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: stack.Name,
		}, broker); err != nil {
			return err
		}
	}

	if broker != nil {
		if !broker.Status.Ready {
			return core.NewPendingError().WithMessage("broker not ready")
		}
		brokerEnvVar, err := brokers.GetBrokerEnvVars(ctx, broker.Status.URI, stack.Name, "payments")
		if err != nil {
			return err
		}

		env = append(env, brokerEnvVar...)
		env = append(env, brokers.GetPublisherEnvVars(stack, broker, "payments", "")...)
	}

	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	err = applications.
		New(payments, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "payments-connectors",
			},
			Spec: appsv1.DeploymentSpec{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						ServiceAccountName: serviceAccountName,
						Containers: []v1.Container{{
							Name:  "connectors",
							Args:  []string{"connectors", "serve"},
							Env:   env,
							Image: image,
							Ports: []v1.ContainerPort{applications.StandardHTTPPort()},
							LivenessProbe: applications.DefaultLiveness("http",
								applications.WithProbePath("/_health")),
						}},
						// Ensure empty
						InitContainers: []v1.Container{},
					},
				},
			},
		}).
		Install(ctx)
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

	caddyfileConfigMap, err := caddy.CreateCaddyfileConfigMap(ctx, stack, "payments", Caddyfile, map[string]any{
		"Debug": stack.Spec.Debug || p.Spec.Debug,
	}, core.WithController[*v1.ConfigMap](ctx.GetScheme(), p))
	if err != nil {
		return err
	}

	env := make([]v1.EnvVar, 0)

	env = append(env, core.GetDevEnvVars(stack, p)...)

	caddyImage, err := registries.GetCaddyImage(ctx, stack, "2.7.6-alpine")
	if err != nil {
		return err
	}

	deploymentTemplate, err := caddy.DeploymentTemplate(ctx, stack, p, caddyfileConfigMap, caddyImage, env)
	if err != nil {
		return err
	}
	// notes(gfyrag): reset init containers in case of upgrading from v1 to v2
	deploymentTemplate.Spec.Template.Spec.InitContainers = make([]v1.Container, 0)

	deploymentTemplate.Name = "payments"

	return applications.
		New(p, deploymentTemplate).
		Install(ctx)
}
