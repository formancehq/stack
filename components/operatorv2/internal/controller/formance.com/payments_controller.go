/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package formance_com

import (
	_ "embed"
	"fmt"
	"github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/v2/internal/resources/databases"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/httpapis"
	"github.com/formancehq/operator/v2/internal/resources/payments"
	. "github.com/formancehq/operator/v2/internal/resources/registries"
	"github.com/formancehq/operator/v2/internal/resources/services"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"github.com/formancehq/operator/v2/internal/resources/streams"
	"github.com/formancehq/operator/v2/internal/resources/topics"
	"github.com/formancehq/search/benthos"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// PaymentsController reconciles a Payments object
type PaymentsController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=payments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=payments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=payments/finalizers,verbs=update

func (r *PaymentsController) Reconcile(ctx Context, payments *v1beta1.Payments) error {

	stack, err := stacks.GetStack(ctx, payments.Spec)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, payments)
	if err != nil {
		return err
	}

	if err := r.createReadDeployment(ctx, stack, payments, database); err != nil {
		return err
	}

	if err := r.createConnectorsDeployment(ctx, stack, payments, database); err != nil {
		return err
	}

	if err := r.createGateway(ctx, stack, payments); err != nil {
		return err
	}

	if err := httpapis.Create(ctx, stack, payments, "payments",
		httpapis.WithRules(
			v1beta1.HTTPAPIRule{
				Path:    "/connectors/webhooks",
				Methods: []string{http.MethodPost},
				Secured: true,
			},
			httpapis.RuleSecured(),
		), httpapis.WithServiceConfiguration(payments.Spec.Service)); err != nil {
		return err
	}

	return nil
}

func (r *PaymentsController) commonEnvVars(ctx Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) ([]corev1.EnvVar, error) {
	env, err := GetCommonServicesEnvVars(ctx, stack, "payments", payments.Spec)
	if err != nil {
		return nil, err
	}
	env = append(env, databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)...)
	env = append(env,
		Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
		Env("CONFIG_ENCRYPTION_KEY", payments.Spec.EncryptionKey),
	)

	return env, nil
}

func (r *PaymentsController) createReadDeployment(ctx Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) error {

	env, _ := r.commonEnvVars(ctx, stack, payments, database)

	image, err := GetImage(ctx, stack, "payments", payments.Spec.Version)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: payments.Spec.Stack,
		Name:      "payments-read",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), payments),
		deployments.WithMatchingLabels("payments-read"),
		deployments.WithContainers(corev1.Container{
			Name:          "api",
			Args:          []string{"api", "serve"},
			Env:           env,
			Image:         image,
			Resources:     GetResourcesWithDefault(payments.Spec.ResourceProperties, ResourceSizeSmall()),
			LivenessProbe: deployments.DefaultLiveness("http", deployments.WithProbePath("/_health")),
			Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
		}),
	)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Namespace: payments.Spec.Stack,
		Name:      "payments-read",
	},
		WithController[*corev1.Service](ctx.GetScheme(), payments),
		services.ConfigureK8SService("payments-read"),
	)
	if err != nil {
		return err
	}

	hasSearch, err := stacks.HasDependency[*v1beta1.Search](ctx, stack.Name)
	if err != nil {
		return err
	}
	if hasSearch {
		if err := streams.LoadFromFileSystem(ctx, benthos.Streams, payments.Spec.Stack, "streams/payments/v0.0.0",
			WithController[*v1beta1.Stream](ctx.GetScheme(), payments),
			WithLabels[*v1beta1.Stream](map[string]string{
				"service": "payments",
			}),
		); err != nil {
			return err
		}
	} else {
		if err := ctx.GetClient().DeleteAllOf(ctx, &v1beta1.Stream{}, client.MatchingLabels{
			"service": "payments",
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *PaymentsController) createConnectorsDeployment(ctx Context, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) error {

	env, _ := r.commonEnvVars(ctx, stack, payments, database)

	topic, err := topics.Find(ctx, stack, "payments")
	if err != nil {
		return err
	}

	if topic != nil {
		if !topic.Status.Ready {
			return fmt.Errorf("topic %s is not yet ready", topic.Name)
		}

		env = append(env, brokerconfigurations.BrokerEnvVars(*topic.Status.Configuration, stack.Name, "payments")...)
		env = append(env, Env("PUBLISHER_TOPIC_MAPPING", "*:"+GetObjectName(stack.Name, "payments")))
	}

	image, err := GetImage(ctx, stack, "payments", payments.Spec.Version)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: payments.Spec.Stack,
		Name:      "payments-connectors",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), payments),
		deployments.WithMatchingLabels("payments-connectors"),
		deployments.WithContainers(corev1.Container{
			Name:      "connectors",
			Args:      []string{"connectors", "serve"},
			Env:       env,
			Image:     image,
			Resources: GetResourcesWithDefault(payments.Spec.ResourceProperties, ResourceSizeSmall()),
			Ports:     []corev1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http",
				deployments.WithProbePath("/_health")),
		}),
		r.setInitContainer(payments, database, image),
	)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Namespace: payments.Spec.Stack,
		Name:      "payments-connectors",
	},
		WithController[*corev1.Service](ctx.GetScheme(), payments),
		services.ConfigureK8SService("payments-connectors"),
	)
	if err != nil {
		return err
	}
	return err
}

func (r *PaymentsController) createGateway(ctx Context, stack *v1beta1.Stack, p *v1beta1.Payments) error {

	caddyfileConfigMap, err := CreateCaddyfileConfigMap(ctx, stack, "payments", payments.Caddyfile, map[string]any{
		"Debug": stack.Spec.Debug || p.Spec.Debug,
	}, WithController[*corev1.ConfigMap](ctx.GetScheme(), p))
	if err != nil {
		return err
	}

	env, err := GetCommonServicesEnvVars(ctx, stack, "payments", p.Spec)
	if err != nil {
		return err
	}

	containerEnv := make([]corev1.EnvVar, 0)
	containerEnv = append(containerEnv, env...)

	mutators := ConfigureCaddy(caddyfileConfigMap, "caddy:2.7.6-alpine", containerEnv, nil)
	mutators = append(mutators,
		WithController[*appsv1.Deployment](ctx.GetScheme(), p),
		deployments.WithMatchingLabels("payments"),
	)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "payments",
	}, mutators...)
	return err
}

func (r *PaymentsController) setInitContainer(payments *v1beta1.Payments, database *v1beta1.Database, image string) func(t *appsv1.Deployment) {
	return func(t *appsv1.Deployment) {
		t.Spec.Template.Spec.InitContainers = []corev1.Container{
			databases.MigrateDatabaseContainer(
				image,
				database.Status.Configuration.DatabaseConfigurationSpec,
				database.Status.Configuration.Database,
				func(m *databases.MigrationConfiguration) {
					m.AdditionalEnv = []corev1.EnvVar{
						Env("CONFIG_ENCRYPTION_KEY", payments.Spec.EncryptionKey),
					}
				},
			),
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PaymentsController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(
			&v1beta1.Database{},
			handler.EnqueueRequestsFromMapFunc(
				databases.Watch[*v1beta1.Payments](mgr, "payments")),
		).
		Watches(
			&v1beta1.Topic{},
			handler.EnqueueRequestsFromMapFunc(
				topics.Watch[*v1beta1.Payments](mgr, "payments")),
		).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Payments](mgr)),
		).
		Watches(
			&v1beta1.RegistriesConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Payments](mgr)),
		).
		Watches(
			&v1beta1.Search{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchDependents[*v1beta1.Ledger](mgr)),
		).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&v1beta1.HTTPAPI{}).
		For(&v1beta1.Payments{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForPayments() *PaymentsController {
	return &PaymentsController{}
}
