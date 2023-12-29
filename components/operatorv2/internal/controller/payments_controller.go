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

package controller

import (
	"context"
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/formancehq/search/benthos"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

//go:embed templates/Caddyfile.payments.gotpl
var paymentsCaddyfile string

// PaymentsController reconciles a Payments object
type PaymentsController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=payments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=payments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=payments/finalizers,verbs=update

func (r *PaymentsController) Reconcile(ctx reconcilers.ContextualManager, payments *v1beta1.Payments) error {

	stack, err := GetStack(ctx, ctx.GetClient(), payments.Spec)
	if err != nil {
		return err
	}

	database, err := CreateDatabase(ctx, ctx.GetClient(), stack, "payments")
	if err != nil {
		return err
	}

	if err := r.createReadDeployment(ctx, stack, payments, database); err != nil {
		return err
	}

	if err := r.createWriteDeployment(ctx, stack, payments, database); err != nil {
		return err
	}

	if err := r.createGateway(ctx, stack, payments); err != nil {
		return err
	}

	if err := CreateHTTPAPI(ctx, ctx.GetClient(), ctx.GetScheme(), stack, payments, "payments"); err != nil {
		return err
	}

	return nil
}

func (r *PaymentsController) commonEnvVars(payments *v1beta1.Payments, database *v1beta1.Database) []corev1.EnvVar {
	env := PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)
	env = append(env,
		Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
		Env("CONFIG_ENCRYPTION_KEY", payments.Spec.EncryptionKey),
	)
	return env
}

func (r *PaymentsController) createReadDeployment(ctx reconcilers.ContextualManager, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) error {

	env := r.commonEnvVars(payments, database)

	_, _, err := CreateOrUpdate[*appsv1.Deployment](ctx, ctx.GetClient(), types.NamespacedName{
		Namespace: payments.Spec.Stack,
		Name:      "payments-read",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), payments),
		WithMatchingLabels("payments-read"),
		WithContainers(corev1.Container{
			Name:      "api",
			Args:      []string{"api", "serve"},
			Env:       env,
			Image:     GetImage("payments", GetVersion(stack, payments.Spec.Version)),
			Resources: GetResourcesWithDefault(payments.Spec.ResourceProperties, ResourceSizeSmall()),
			Ports:     []corev1.ContainerPort{StandardHTTPPort()},
		}),
	)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*corev1.Service](ctx, ctx.GetClient(), types.NamespacedName{
		Namespace: payments.Spec.Stack,
		Name:      "payments-read",
	},
		WithController[*corev1.Service](ctx.GetScheme(), payments),
		ConfigureK8SService("payments-read"),
	)
	if err != nil {
		return err
	}

	if err := LoadStreams(ctx, ctx.GetClient(), benthos.Streams, payments.Spec.Stack, "streams/payments/v0.0.0"); err != nil {
		return err
	}

	return nil
}

func (r *PaymentsController) createWriteDeployment(ctx reconcilers.ContextualManager, stack *v1beta1.Stack, payments *v1beta1.Payments, database *v1beta1.Database) error {

	env := r.commonEnvVars(payments, database)

	topicExists, err := TopicExists(ctx, ctx.GetClient(), stack, "payments")
	if err != nil {
		return err
	}

	if topicExists {
		brokerEnvVars, err := GetBrokerEnvVars(ctx, ctx.GetClient(), stack.Name, "payments")
		if err != nil {
			return err
		}
		env = append(env, brokerEnvVars...)
		env = append(env, Env("PUBLISHER_TOPIC_MAPPING", "*:"+GetObjectName(stack.Name, "payments")))
	}

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, ctx.GetClient(), types.NamespacedName{
		Namespace: payments.Spec.Stack,
		Name:      "payments-connectors",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), payments),
		WithMatchingLabels("payments-connectors"),
		WithContainers(corev1.Container{
			Name:      "connectors",
			Args:      []string{"connectors", "serve"},
			Env:       env,
			Image:     GetImage("payments", GetVersion(stack, payments.Spec.Version)),
			Resources: GetResourcesWithDefault(payments.Spec.ResourceProperties, ResourceSizeSmall()),
			Ports:     []corev1.ContainerPort{StandardHTTPPort()},
		}),
		r.setInitContainer(payments, database, GetVersion(stack, payments.Spec.Version)),
	)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*corev1.Service](ctx, ctx.GetClient(), types.NamespacedName{
		Namespace: payments.Spec.Stack,
		Name:      "payments-connectors",
	},
		WithController[*corev1.Service](ctx.GetScheme(), payments),
		ConfigureK8SService("payments-connectors"),
	)
	if err != nil {
		return err
	}
	return err
}

func (r *PaymentsController) createGateway(ctx reconcilers.ContextualManager, stack *v1beta1.Stack, payments *v1beta1.Payments) error {

	caddyfileConfigMap, err := CreateCaddyfileConfigMap(ctx, ctx.GetClient(), stack, "payments", paymentsCaddyfile, map[string]any{
		"Debug": stack.Spec.Debug || payments.Spec.Debug,
	}, WithController[*corev1.ConfigMap](ctx.GetScheme(), payments))
	if err != nil {
		return err
	}

	env, err := GetCommonServicesEnvVars(ctx, ctx.GetClient(), stack, "payments", payments.Spec)
	if err != nil {
		return err
	}

	containerEnv := make([]corev1.EnvVar, 0)
	containerEnv = append(containerEnv, env...)

	mutators := ConfigureCaddy(caddyfileConfigMap, "caddy:2.7.6-alpine", containerEnv, nil)
	mutators = append(mutators,
		WithController[*appsv1.Deployment](ctx.GetScheme(), payments),
		WithMatchingLabels("payments"),
	)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, ctx.GetClient(), types.NamespacedName{
		Namespace: stack.Name,
		Name:      "payments",
	}, mutators...)
	return err
}

func (r *PaymentsController) setInitContainer(payments *v1beta1.Payments, database *v1beta1.Database, version string) func(t *appsv1.Deployment) {
	return func(t *appsv1.Deployment) {
		t.Spec.Template.Spec.InitContainers = []corev1.Container{
			MigrateDatabaseContainer(
				"payments",
				database.Status.Configuration.DatabaseConfigurationSpec,
				database.Status.Configuration.Database,
				version,
				func(m *MigrationConfiguration) {
					m.AdditionalEnv = []corev1.EnvVar{
						Env("CONFIG_ENCRYPTION_KEY", payments.Spec.EncryptionKey),
					}
				},
			),
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PaymentsController) SetupWithManager(mgr reconcilers.Manager) (*builder.Builder, error) {

	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.Payments{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.Payments).Spec.Stack}
	}); err != nil {
		return nil, err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Payments{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForPayments() *PaymentsController {
	return &PaymentsController{}
}
