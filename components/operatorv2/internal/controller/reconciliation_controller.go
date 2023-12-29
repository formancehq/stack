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
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/formancehq/operator/v2/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ReconciliationController reconciles a Reconciliation object
type ReconciliationController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=reconciliations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=reconciliations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=reconciliations/finalizers,verbs=update

func (r *ReconciliationController) Reconcile(ctx reconcilers.ContextualManager, reconciliation *v1beta1.Reconciliation) error {
	stack, err := GetStack(ctx, ctx.GetClient(), reconciliation.Spec)
	if err != nil {
		return err
	}

	database, err := CreateDatabase(ctx, ctx.GetClient(), stack, "reconciliation")
	if err != nil {
		return err
	}

	authClient, err := CreateAuthClient(ctx, ctx.GetClient(), ctx.GetScheme(), stack, reconciliation, "reconciliation", func(spec *v1beta1.AuthClientSpec) {
		spec.Scopes = []string{"ledger:read", "payments:read"}
	})
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, reconciliation, database, authClient); err != nil {
		return err
	}

	if err := CreateHTTPAPI(ctx, ctx.GetClient(), ctx.GetScheme(), stack, reconciliation, "reconciliation"); err != nil {
		return err
	}

	return nil
}

func (r *ReconciliationController) createDeployment(ctx reconcilers.ContextualManager, stack *v1beta1.Stack,
	reconciliation *v1beta1.Reconciliation, database *v1beta1.Database, authClient *v1beta1.AuthClient) error {
	env, err := GetCommonServicesEnvVars(ctx, ctx.GetClient(), stack, "reconciliation", reconciliation.Spec)
	if err != nil {
		return err
	}

	env = append(env,
		PostgresEnvVars(
			database.Status.Configuration.DatabaseConfigurationSpec,
			GetObjectName(stack.Name, "reconciliation"),
		)...,
	)
	env = append(env,
		Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
	)
	env = append(env, GetAuthClientEnvVars(authClient)...)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, ctx.GetClient(),
		GetNamespacedResourceName(stack.Name, "reconciliation"),
		func(t *appsv1.Deployment) {
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:      "reconciliation",
				Env:       env,
				Image:     GetImage("reconciliation", GetVersion(stack, reconciliation.Spec.Version)),
				Resources: GetResourcesWithDefault(reconciliation.Spec.ResourceProperties, ResourceSizeSmall()),
				Ports:     []corev1.ContainerPort{StandardHTTPPort()},
			}}
		},
		WithMatchingLabels("reconciliation"),
		WithController[*appsv1.Deployment](ctx.GetScheme(), reconciliation),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReconciliationController) SetupWithManager(mgr reconcilers.Manager) (*builder.Builder, error) {

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Reconciliation{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForReconciliation() *ReconciliationController {
	return &ReconciliationController{}
}
