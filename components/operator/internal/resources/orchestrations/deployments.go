package orchestrations

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/internal/resources/brokers"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/applications"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/resourcereferences"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
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

	env := make([]corev1.EnvVar, 0)
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

	env = append(env, gatewayEnv...)
	env = append(env, GetDevEnvVars(stack, orchestration)...)
	env = append(env, postgresEnvVar...)

	temporalURI, err := settings.RequireURL(ctx, stack.Name, "temporal", "dsn")
	if err != nil {
		return err
	}

	if err := validateTemporalURI(temporalURI); err != nil {
		return err
	}

	var temporalSecretResourceReference *v1beta1.ResourceReference
	if secret := temporalURI.Query().Get("secret"); secret != "" {
		temporalSecretResourceReference, err = resourcereferences.Create(ctx, orchestration, "temporal", secret, &corev1.Secret{})
	} else {
		err = resourcereferences.Delete(ctx, orchestration, "temporal")
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

	broker := &v1beta1.Broker{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: stack.Name,
	}, broker); err != nil {
		return err
	}

	brokerEnvVars, err := brokers.GetBrokerEnvVars(ctx, broker.Status.URI, stack.Name, "orchestration")
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	env = append(env, brokerEnvVars...)
	env = append(env, brokers.GetPublisherEnvVars(stack, broker, "orchestration", "")...)

	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	annotations := map[string]string{}
	if temporalSecretResourceReference != nil {
		annotations["database-secret-hash"] = temporalSecretResourceReference.Status.Hash
	}

	maxParallelActivities, err := settings.GetIntOrDefault(ctx, stack.Name, 10, "orchestration", "max-parallel-activities")
	if err != nil {
		return err
	}
	env = append(env, Env("TEMPORAL_MAX_PARALLEL_ACTIVITIES", fmt.Sprint(maxParallelActivities)))

	tpl := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "orchestration",
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Containers: []corev1.Container{{
						Name:          "api",
						Env:           env,
						Image:         image,
						Ports:         []corev1.ContainerPort{applications.StandardHTTPPort()},
						LivenessProbe: applications.DefaultLiveness("http"),
					}},
				},
			},
		},
	}

	return applications.
		New(orchestration, tpl).
		IsEE().
		Install(ctx)
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
