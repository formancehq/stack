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
	"github.com/formancehq/operator/v2/internal/brokerconfigurations"
	common "github.com/formancehq/operator/v2/internal/common"
	databases "github.com/formancehq/operator/v2/internal/databases"
	"github.com/formancehq/operator/v2/internal/deployments"
	"github.com/formancehq/operator/v2/internal/httpapis"
	"github.com/formancehq/operator/v2/internal/ledgers"
	payments2 "github.com/formancehq/operator/v2/internal/payments"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/formancehq/operator/v2/internal/topicqueries"
	. "github.com/formancehq/operator/v2/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"

	v1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// WebhooksController reconciles a Webhooks object
type WebhooksController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=webhooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=webhooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=webhooks/finalizers,verbs=update

func (r *WebhooksController) Reconcile(ctx reconcilers.Context, webhooks *v1beta1.Webhooks) error {
	stack, err := common.GetStack(ctx, webhooks.Spec)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, stack, "webhooks")
	if err != nil {
		return err
	}

	if err := r.handleTopics(ctx, stack); err != nil {
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
func (r *WebhooksController) handleTopics(ctx reconcilers.Context, stack *v1beta1.Stack) error {
	availableServices := make([]string, 0)
	ledger, err := ledgers.GetIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	if ledger != nil {
		availableServices = append(availableServices, "ledger")
	}
	payments, err := payments2.GetIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	if payments != nil {
		availableServices = append(availableServices, "payments")
	}

	for _, service := range availableServices {
		if err := topicqueries.Create(ctx, stack, service, "webhooks"); err != nil {
			return err
		}
	}

	return nil
}

func (r *WebhooksController) createDeployment(ctx reconcilers.Context, stack *v1beta1.Stack, webhooks *v1beta1.Webhooks,
	database *v1beta1.Database) error {

	brokerConfiguration, err := brokerconfigurations.Require(ctx, stack.Name)
	if err != nil {
		return err
	}

	env := databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)
	env = append(env, brokerconfigurations.BrokerEnvVars(*brokerConfiguration, "webhooks")...)
	env = append(env, Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))
	env = append(env, Env("KAFKA_TOPICS", strings.Join([]string{
		common.GetObjectName(stack.Name, "ledger"),
		common.GetObjectName(stack.Name, "payments"),
	}, " ")))
	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: webhooks.Spec.Stack,
		Name:      "webhooks",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), webhooks),
		deployments.WithMatchingLabels("webhooks"),
		deployments.WithContainers(corev1.Container{
			Name:      "api",
			Env:       env,
			Image:     common.GetImage("webhooks", common.GetVersion(stack, webhooks.Spec.Version)),
			Resources: GetResourcesWithDefault(webhooks.Spec.ResourceProperties, ResourceSizeSmall()),
			Ports:     []corev1.ContainerPort{common.StandardHTTPPort()},
		}),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebhooksController) SetupWithManager(mgr reconcilers.Manager) (*builder.Builder, error) {

	return ctrl.NewControllerManagedBy(mgr).
		//TODO: Watch services able to trigger events
		For(&v1beta1.Webhooks{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForWebhooks() *WebhooksController {
	return &WebhooksController{}
}
