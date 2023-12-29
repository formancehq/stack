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
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/authclients"
	"github.com/formancehq/operator/v2/internal/resources/auths"
	"github.com/formancehq/operator/v2/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/v2/internal/resources/databases"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/httpapis"
	"github.com/formancehq/operator/v2/internal/resources/ledgers"
	"github.com/formancehq/operator/v2/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/v2/internal/resources/payments"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"github.com/formancehq/operator/v2/internal/resources/topicqueries"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	formancev1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// OrchestrationController reconciles a Orchestration object
type OrchestrationController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=orchestrations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/finalizers,verbs=update

func (r *OrchestrationController) Reconcile(ctx core.Context, orchestration *v1beta1.Orchestration) error {

	stack, err := stacks.GetStack(ctx, orchestration.Spec)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, stack, "orchestration")
	if err != nil {
		return err
	}

	authClient, err := r.handleAuthClient(ctx, stack, orchestration)
	if err != nil {
		return err
	}

	if err := r.handleTopics(ctx, stack); err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, orchestration, database, authClient); err != nil {
		return err
	}

	if err := httpapis.Create(ctx, stack, orchestration, "orchestration"); err != nil {
		return err
	}

	return nil
}

func (r *OrchestrationController) handleAuthClient(ctx core.Context, stack *formancev1beta1.Stack, orchestration *formancev1beta1.Orchestration) (*formancev1beta1.AuthClient, error) {

	auth, err := auths.GetIfEnabled(ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, errors.New("requires auth service")
	}

	return authclients.Create(ctx, stack, orchestration, "orchestration",
		func(spec *formancev1beta1.AuthClientSpec) {
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

func (r *OrchestrationController) handleTopics(ctx core.Context, stack *formancev1beta1.Stack) error {
	availableServices := make([]string, 0)
	ledger, err := ledgers.GetIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	if ledger != nil {
		availableServices = append(availableServices, "ledger")
	}
	payments, err := payments.GetIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	if payments != nil {
		availableServices = append(availableServices, "payments")
	}

	for _, service := range availableServices {
		if err := topicqueries.Create(ctx, stack, service, "orchestration"); err != nil {
			return err
		}
	}

	return nil
}

func (r *OrchestrationController) createDeployment(ctx core.Context, stack *v1beta1.Stack,
	orchestration *v1beta1.Orchestration, database *v1beta1.Database, client *formancev1beta1.AuthClient) error {
	env := databases.PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)
	env = append(env,
		core.Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		core.Env("TEMPORAL_TASK_QUEUE", stack.Name),
		core.Env("TEMPORAL_ADDRESS", orchestration.Spec.Temporal.Address),
		core.Env("TEMPORAL_NAMESPACE", orchestration.Spec.Temporal.Namespace),
	)

	env = append(env, authclients.GetEnvVars(client)...)

	if orchestration.Spec.Temporal.TLS.SecretName == "" {
		env = append(env,
			core.Env("TEMPORAL_SSL_CLIENT_KEY", orchestration.Spec.Temporal.TLS.Key),
			core.Env("TEMPORAL_SSL_CLIENT_CERT", orchestration.Spec.Temporal.TLS.CRT),
		)
	} else {
		env = append(env,
			core.EnvFromSecret("TEMPORAL_SSL_CLIENT_KEY", orchestration.Spec.Temporal.TLS.SecretName, "tls.key"),
			core.EnvFromSecret("TEMPORAL_SSL_CLIENT_CERT", orchestration.Spec.Temporal.TLS.SecretName, "tls.crt"),
		)
	}

	brokerEnvVars, err := brokerconfigurations.GetEnvVars(ctx, stack.Name, "orchestration")
	if err != nil && !errors.Is(err, stacks.ErrNotFound) {
		return err
	}
	env = append(env, brokerEnvVars...)

	_, _, err = core.CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: orchestration.Spec.Stack,
		Name:      "orchestration",
	},
		core.WithController[*appsv1.Deployment](ctx.GetScheme(), orchestration),
		deployments.WithMatchingLabels("orchestration"),
		deployments.WithContainers(corev1.Container{
			Name:      "api",
			Env:       env,
			Image:     core.GetImage("orchestration", core.GetVersion(stack, orchestration.Spec.Version)),
			Resources: core.GetResourcesWithDefault(orchestration.Spec.ResourceProperties, core.ResourceSizeSmall()),
			Ports:     []corev1.ContainerPort{deployments.StandardHTTPPort()},
		}),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrchestrationController) SetupWithManager(mgr core.Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(
			&v1beta1.Database{},
			handler.EnqueueRequestsFromMapFunc(databases.Watch[*v1beta1.OrchestrationList, *v1beta1.Orchestration](mgr, "ledger")),
		).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(
				opentelemetryconfigurations.Watch[*v1beta1.OrchestrationList, *v1beta1.Orchestration](mgr),
			),
		).
		// todo: Watch broker configuration
		For(&formancev1beta1.Orchestration{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForOrchestration() *OrchestrationController {
	return &OrchestrationController{}
}
