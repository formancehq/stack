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
	"fmt"
	v1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	deployments "github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/httpapis"
	. "github.com/formancehq/operator/v2/internal/resources/registries"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// SearchController reconciles a Search object
type SearchController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=searches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=searches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=searches/finalizers,verbs=update

func (r *SearchController) Reconcile(ctx Context, search *v1beta1.Search) error {

	stack, err := stacks.GetStack(ctx, search.Spec)
	if err != nil {
		return err
	}

	elasticSearchConfiguration, err := stacks.Require[*v1beta1.ElasticSearchConfiguration](ctx, search.Spec.Stack)
	if err != nil {
		return err
	}

	env, err := GetCommonServicesEnvVars(ctx, stack, "search", search.Spec)
	if err != nil {
		return err
	}

	env = append(env,
		Env("OPEN_SEARCH_SERVICE", fmt.Sprintf("%s:%d", elasticSearchConfiguration.Spec.Host, elasticSearchConfiguration.Spec.Port)),
		Env("OPEN_SEARCH_SCHEME", elasticSearchConfiguration.Spec.Scheme),
		Env("ES_INDICES", "stacks"),
	)
	if elasticSearchConfiguration.Spec.BasicAuth != nil {
		if elasticSearchConfiguration.Spec.BasicAuth.SecretName == "" {
			env = append(env,
				Env("OPEN_SEARCH_USERNAME", elasticSearchConfiguration.Spec.BasicAuth.Username),
				Env("OPEN_SEARCH_PASSWORD", elasticSearchConfiguration.Spec.BasicAuth.Password),
			)
		} else {
			env = append(env,
				EnvFromSecret("OPEN_SEARCH_USERNAME", elasticSearchConfiguration.Spec.BasicAuth.SecretName, "username"),
				EnvFromSecret("OPEN_SEARCH_PASSWORD", elasticSearchConfiguration.Spec.BasicAuth.SecretName, "password"),
			)
		}
	}

	image, err := GetImage(ctx, stack, "search", search.Spec.Version)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*v1beta1.StreamProcessor](ctx, types.NamespacedName{
		Name: GetObjectName(stack.Name, "stream-processor"),
	},
		WithController[*v1beta1.StreamProcessor](ctx.GetScheme(), search),
		func(t *v1beta1.StreamProcessor) {
			t.Spec.Stack = stack.Name
			t.Spec.Batching = search.Spec.Batching
			t.Spec.DevProperties = search.Spec.DevProperties
			t.Spec.ResourceProperties = search.Spec.CommonServiceProperties.ResourceProperties
			t.Spec.InitContainers = []corev1.Container{{
				Name:  "init-mapping",
				Image: image,
				Args:  []string{"init-mapping"},
				Env:   env,
			}}
		},
	)
	if err != nil {
		return err
	}

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "search",
	},
		WithController[*appsv1.Deployment](ctx.GetScheme(), search),
		deployments.WithMatchingLabels("search"),
		deployments.WithContainers(corev1.Container{
			Name:            "search",
			Image:           image,
			Ports:           []corev1.ContainerPort{deployments.StandardHTTPPort()},
			Env:             env,
			Resources:       GetResourcesWithDefault(search.Spec.ResourceProperties, ResourceSizeSmall()),
			ImagePullPolicy: GetPullPolicy(image),
			LivenessProbe:   deployments.DefaultLiveness("http"),
		}),
	)

	if err := httpapis.Create(ctx, stack, search, "search"); err != nil {
		return err
	}

	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SearchController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Search{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Search](mgr)),
		).
		Watches(
			&v1beta1.ElasticSearchConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Search](mgr)),
		).
		Watches(
			&v1beta1.Registries{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Search](mgr)),
		).
		Owns(&v1beta1.StreamProcessor{}).
		Owns(&v1beta1.HTTPAPI{}).
		Owns(&appsv1.Deployment{}), nil
}

func ForSearch() *SearchController {
	return &SearchController{}
}