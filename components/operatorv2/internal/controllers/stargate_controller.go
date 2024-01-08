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
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/registries"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// StargateController reconciles a Stargate object
type StargateController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=stargates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=stargates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=stargates/finalizers,verbs=update

func (r *StargateController) Reconcile(ctx Context, stargate *v1beta1.Stargate) error {

	stack, err := stacks.GetStack(ctx, stargate.Spec)
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, stargate); err != nil {
		return err
	}

	return nil
}

func (r *StargateController) createDeployment(ctx Context, stack *v1beta1.Stack, stargate *v1beta1.Stargate) error {

	env, err := GetCommonServicesEnvVars(ctx, stack, "wallets", stargate.Spec)
	if err != nil {
		return err
	}
	env = append(env,
		Env("ORGANIZATION_ID", stargate.Spec.OrganizationID),
		Env("STACK_ID", stargate.Spec.StackID),
		Env("STARGATE_SERVER_URL", stargate.Spec.ServerURL),
		Env("GATEWAY_URL", "http://gateway:8080"),
		Env("STARGATE_AUTH_CLIENT_ID", stargate.Spec.Auth.ClientID),
		Env("STARGATE_AUTH_CLIENT_SECRET", stargate.Spec.Auth.ClientSecret),
		Env("STARGATE_AUTH_ISSUER_URL", stargate.Spec.Auth.Issuer),
	)

	image, err := registries.GetImage(ctx, stack, "stargate", stargate.Spec.Version)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx,
		GetNamespacedResourceName(stack.Name, "stargate"),
		func(t *appsv1.Deployment) {
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:          "stargate",
				Env:           env,
				Image:         image,
				Resources:     GetResourcesWithDefault(stargate.Spec.ResourceProperties, ResourceSizeSmall()),
				Ports:         []corev1.ContainerPort{deployments.StandardHTTPPort()},
				LivenessProbe: deployments.DefaultLiveness("http"),
			}}
		},
		deployments.WithMatchingLabels("stargate"),
		WithController[*appsv1.Deployment](ctx.GetScheme(), stargate),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *StargateController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Stargate](mgr)),
		).
		Watches(
			&v1beta1.Registries{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Stargate](mgr)),
		).
		Owns(&appsv1.Deployment{}).
		For(&v1beta1.Stargate{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForStargate() *StargateController {
	return &StargateController{}
}
