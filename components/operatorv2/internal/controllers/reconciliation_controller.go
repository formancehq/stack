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
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/authclients"
	"github.com/formancehq/operator/v2/internal/resources/databases"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/httpapis"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/formancehq/operator/v2/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ReconciliationController reconciles a Reconciliation object
type ReconciliationController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=reconciliations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=reconciliations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=reconciliations/finalizers,verbs=update

func (r *ReconciliationController) Reconcile(ctx Context, reconciliation *v1beta1.Reconciliation) error {
	stack, err := stacks.GetStack(ctx, reconciliation.Spec)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, reconciliation)
	if err != nil {
		return err
	}

	authClient, err := authclients.Create(ctx, stack, reconciliation, "reconciliation",
		authclients.WithScopes("ledger:read", "payments:read"))
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, reconciliation, database, authClient); err != nil {
		return err
	}

	if err := httpapis.Create(ctx, stack, reconciliation, "reconciliation"); err != nil {
		return err
	}

	return nil
}

func (r *ReconciliationController) createDeployment(ctx Context, stack *v1beta1.Stack,
	reconciliation *v1beta1.Reconciliation, database *v1beta1.Database, authClient *v1beta1.AuthClient) error {
	env, err := GetCommonServicesEnvVars(ctx, stack, "reconciliation", reconciliation.Spec)
	if err != nil {
		return err
	}

	env = append(env,
		databases.PostgresEnvVars(
			database.Status.Configuration.DatabaseConfigurationSpec,
			GetObjectName(stack.Name, "reconciliation"),
		)...,
	)
	env = append(env,
		Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
	)
	env = append(env, authclients.GetEnvVars(authClient)...)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx,
		GetNamespacedResourceName(stack.Name, "reconciliation"),
		func(t *appsv1.Deployment) {
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:          "reconciliation",
				Env:           env,
				Image:         GetImage("reconciliation", GetVersion(stack, reconciliation.Spec.Version)),
				Resources:     GetResourcesWithDefault(reconciliation.Spec.ResourceProperties, ResourceSizeSmall()),
				Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
				LivenessProbe: deployments.DefaultLiveness("http"),
			}}
		},
		deployments.WithMatchingLabels("reconciliation"),
		WithController[*appsv1.Deployment](ctx.GetScheme(), reconciliation),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReconciliationController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(
			&v1beta1.Ledger{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchDependents(mgr, &v1beta1.ReconciliationList{})),
		).
		Watches(
			&v1beta1.Payments{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchDependents(mgr, &v1beta1.ReconciliationList{})),
		).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(
				Watch(mgr, &v1beta1.ReconciliationList{}),
			),
		).
		Owns(&v1beta1.Database{}).
		Owns(&appsv1.Deployment{}).
		Owns(&v1beta1.AuthClient{}).
		Owns(&v1beta1.HTTPAPI{}).
		For(&v1beta1.Reconciliation{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForReconciliation() *ReconciliationController {
	return &ReconciliationController{}
}
