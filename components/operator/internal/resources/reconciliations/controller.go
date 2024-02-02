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

package reconciliations

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/gatewayhttpapis"
	appsv1 "k8s.io/api/apps/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=reconciliations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=reconciliations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=reconciliations/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, reconciliation *v1beta1.Reconciliation, version string) error {
	database, err := databases.Create(ctx, reconciliation)
	if err != nil {
		return err
	}

	authClient, err := authclients.Create(ctx, stack, reconciliation, "reconciliation",
		authclients.WithScopes("ledger:read", "payments:read"))
	if err != nil {
		return err
	}

	if database.Status.Ready {
		if err := createDeployment(ctx, stack, reconciliation, database, authClient, version); err != nil {
			return err
		}
	}

	if err := gatewayhttpapis.Create(ctx, reconciliation); err != nil {
		return err
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithOwn[*v1beta1.Reconciliation](&v1beta1.Database{}),
			WithOwn[*v1beta1.Reconciliation](&appsv1.Deployment{}),
			WithOwn[*v1beta1.Reconciliation](&v1beta1.AuthClient{}),
			WithOwn[*v1beta1.Reconciliation](&v1beta1.GatewayHTTPAPI{}),
			WithWatchSettings[*v1beta1.Reconciliation](),
			WithWatchDependency[*v1beta1.Reconciliation](&v1beta1.Ledger{}),
			WithWatchDependency[*v1beta1.Reconciliation](&v1beta1.Payments{}),
		),
	)
}
