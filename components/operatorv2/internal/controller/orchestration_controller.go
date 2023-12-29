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
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	formancev1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// OrchestrationController reconciles a Orchestration object
type OrchestrationController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=orchestrations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/finalizers,verbs=update

func (r *OrchestrationController) Reconcile(ctx reconcilers.ContextualManager, orchestration *v1beta1.Orchestration) error {

	stack, err := GetStack(ctx, ctx.GetClient(), orchestration.Spec)
	if err != nil {
		return err
	}

	database, err := CreateDatabase(ctx, ctx.GetClient(), stack, "orchestration")
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

	if err := CreateHTTPAPI(ctx, ctx.GetClient(), ctx.GetScheme(), stack, orchestration, "orchestration"); err != nil {
		return err
	}

	return nil
}

func (r *OrchestrationController) handleAuthClient(ctx reconcilers.ContextualManager, stack *formancev1beta1.Stack, orchestration *formancev1beta1.Orchestration) (*formancev1beta1.AuthClient, error) {

	auth, err := GetAuthIfEnabled(ctx, ctx.GetClient(), stack.Name)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, errors.New("requires auth service")
	}

	return CreateAuthClient(ctx, ctx.GetClient(), ctx.GetScheme(), stack, orchestration, "orchestration",
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

func (r *OrchestrationController) handleTopics(ctx reconcilers.ContextualManager, stack *formancev1beta1.Stack) error {
	availableServices := make([]string, 0)
	ledger, err := GetLedgerIfEnabled(ctx, ctx.GetClient(), stack.Name)
	if err != nil {
		return err
	}
	if ledger != nil {
		availableServices = append(availableServices, "ledger")
	}
	payments, err := GetPaymentsIfEnabled(ctx, ctx.GetClient(), stack.Name)
	if err != nil {
		return err
	}
	if payments != nil {
		availableServices = append(availableServices, "payments")
	}

	for _, service := range availableServices {
		if err := CreateTopicQuery(ctx, ctx.GetClient(), stack, service, "orchestration"); err != nil {
			return err
		}
	}

	return nil
}

func (r *OrchestrationController) createDeployment(ctx reconcilers.ContextualManager, stack *v1beta1.Stack,
	orchestration *v1beta1.Orchestration, database *v1beta1.Database, client *formancev1beta1.AuthClient) error {
	env := PostgresEnvVars(database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)
	env = append(env,
		Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		Env("TEMPORAL_TASK_QUEUE", stack.Name),
		Env("TEMPORAL_ADDRESS", orchestration.Spec.Temporal.Address),
		Env("TEMPORAL_NAMESPACE", orchestration.Spec.Temporal.Namespace),
	)

	env = append(env, GetAuthClientEnvVars(client)...)

	if orchestration.Spec.Temporal.TLS.SecretName == "" {
		env = append(env,
			Env("TEMPORAL_SSL_CLIENT_KEY", orchestration.Spec.Temporal.TLS.Key),
			Env("TEMPORAL_SSL_CLIENT_CERT", orchestration.Spec.Temporal.TLS.CRT),
		)
	} else {
		env = append(env,
			EnvFromSecret("TEMPORAL_SSL_CLIENT_KEY", orchestration.Spec.Temporal.TLS.SecretName, "tls.key"),
			EnvFromSecret("TEMPORAL_SSL_CLIENT_CERT", orchestration.Spec.Temporal.TLS.SecretName, "tls.crt"),
		)
	}

	brokerEnvVars, err := GetBrokerEnvVars(ctx, ctx.GetClient(), stack.Name, "orchestration")
	if err != nil && !errors.Is(err, ErrNoConfigurationFound) {
		return err
	}
	env = append(env, brokerEnvVars...)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, ctx.GetClient(), types.NamespacedName{
		Namespace: orchestration.Spec.Stack,
		Name:      "orchestration",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), orchestration),
		WithMatchingLabels("orchestration"),
		WithContainers(corev1.Container{
			Name:      "api",
			Env:       env,
			Image:     GetImage("orchestration", GetVersion(stack, orchestration.Spec.Version)),
			Resources: GetResourcesWithDefault(orchestration.Spec.ResourceProperties, ResourceSizeSmall()),
			Ports:     []corev1.ContainerPort{StandardHTTPPort()},
		}),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrchestrationController) SetupWithManager(mgr reconcilers.Manager) (*builder.Builder, error) {

	return ctrl.NewControllerManagedBy(mgr).
		// todo: Watch broker configuration
		For(&formancev1beta1.Orchestration{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForOrchestration() *OrchestrationController {
	return &OrchestrationController{}
}
