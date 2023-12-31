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

package controllers

import (
	v1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/v2/internal/resources/databases"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/httpapis"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"github.com/formancehq/operator/v2/internal/resources/topicqueries"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

// WebhooksController reconciles a Webhooks object
type WebhooksController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=webhooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=webhooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=webhooks/finalizers,verbs=update

func (r *WebhooksController) Reconcile(ctx Context, webhooks *v1beta1.Webhooks) error {
	stack, err := stacks.GetStack(ctx, webhooks.Spec)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, webhooks)
	if err != nil {
		return err
	}

	if err := r.handleTopics(ctx, stack, webhooks); err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, webhooks, database); err != nil {
		return err
	}

	if err := httpapis.Create(ctx, stack, webhooks, "webhooks"); err != nil {
		return err
	}

	return nil
}

// TODO: Search a way to automatically list all services able to push events as this code is duplicated for orchestration
func (r *WebhooksController) handleTopics(ctx Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks) error {
	return ForEachEventPublisher(ctx, stack.Name, func(object client.Object) error {
		return topicqueries.Create(ctx, stack, strings.ToLower(object.GetObjectKind().GroupVersionKind().Kind), webhooks)
	})
}

func (r *WebhooksController) createDeployment(ctx Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks,
	database *v1beta1.Database) error {

	brokerConfiguration, err := brokerconfigurations.Require(ctx, stack.Name)
	if err != nil {
		return err
	}

	env := databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)
	env = append(env, brokerconfigurations.BrokerEnvVars(brokerConfiguration.Spec, "webhooks")...)
	env = append(env, Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))
	env = append(env, Env("KAFKA_TOPICS", strings.Join([]string{
		GetObjectName(stack.Name, "ledger"),
		GetObjectName(stack.Name, "payments"),
	}, " ")))
	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: webhooks.Spec.Stack,
		Name:      "webhooks",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), webhooks),
		deployments.WithMatchingLabels("webhooks"),
		deployments.WithContainers(corev1.Container{
			Name:          "api",
			Env:           env,
			Image:         GetImage("webhooks", GetVersion(stack, webhooks.Spec.Version)),
			Resources:     GetResourcesWithDefault(webhooks.Spec.ResourceProperties, ResourceSizeSmall()),
			Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
			LivenessProbe: deployments.DefaultLiveness("http"),
		}),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebhooksController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		Owns(&appsv1.Deployment{}).
		Owns(&v1beta1.HTTPAPI{}).
		Watches(
			&v1beta1.Ledger{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents(mgr, &v1beta1.WebhooksList{})),
		).
		Watches(
			&v1beta1.Payments{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents(mgr, &v1beta1.WebhooksList{})),
		).
		Watches(
			&v1beta1.Database{},
			handler.EnqueueRequestsFromMapFunc(
				databases.Watch(mgr, "webhooks", &v1beta1.WebhooksList{})),
		).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(
				Watch(mgr, &v1beta1.WebhooksList{}),
			),
		).
		For(&v1beta1.Webhooks{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForWebhooks() *WebhooksController {
	return &WebhooksController{}
}
