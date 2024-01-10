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
	"github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/authclients"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/httpapis"
	"github.com/formancehq/operator/v2/internal/resources/registries"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	ctrl "sigs.k8s.io/controller-runtime"
)

// WalletsController reconciles a Wallets object
type WalletsController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=wallets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=wallets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=wallets/finalizers,verbs=update

func (r *WalletsController) Reconcile(ctx Context, wallets *v1beta1.Wallets) error {

	stack, err := stacks.GetStack(ctx, wallets.Spec)
	if err != nil {
		return err
	}

	authClient, err := authclients.Create(ctx, stack, wallets, "wallets", func(spec *v1beta1.AuthClientSpec) {
		spec.Scopes = []string{"ledger:read", "ledger:write"}
	})
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, wallets, authClient); err != nil {
		return err
	}

	if err := httpapis.Create(ctx, stack, wallets, "wallets",
		httpapis.WithServiceConfiguration(wallets.Spec.Service)); err != nil {
		return err
	}

	return nil
}

func (r *WalletsController) createDeployment(ctx Context, stack *v1beta1.Stack, wallets *v1beta1.Wallets, authClient *v1beta1.AuthClient) error {
	env, err := GetCommonServicesEnvVars(ctx, stack, "wallets", wallets.Spec)
	if err != nil {
		return err
	}
	env = append(env, authclients.GetEnvVars(authClient)...)

	image, err := registries.GetImage(ctx, stack, "wallets", wallets.Spec.Version)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx,
		GetNamespacedResourceName(stack.Name, "wallets"),
		func(t *appsv1.Deployment) {
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:          "wallets",
				Args:          []string{"serve"},
				Env:           env,
				Image:         image,
				Resources:     GetResourcesWithDefault(wallets.Spec.ResourceProperties, ResourceSizeSmall()),
				Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
				LivenessProbe: deployments.DefaultLiveness("http"),
			}}
		},
		deployments.WithMatchingLabels("wallets"),
		WithController[*appsv1.Deployment](ctx.GetScheme(), wallets),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *WalletsController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Wallets](mgr)),
		).
		Watches(
			&v1beta1.RegistriesConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Wallets](mgr)),
		).
		Owns(&v1beta1.AuthClient{}).
		Owns(&appsv1.Deployment{}).
		Owns(&v1beta1.HTTPAPI{}).
		For(&v1beta1.Wallets{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForWallets() *WalletsController {
	return &WalletsController{}
}
