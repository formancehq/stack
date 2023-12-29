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
	common "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/deployments"
	"github.com/formancehq/operator/v2/internal/elasticsearchconfigurations"
	"github.com/formancehq/operator/v2/internal/httpapis"
	"github.com/formancehq/operator/v2/internal/stacks"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// SearchController reconciles a Search object
type SearchController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=searches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=searches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=searches/finalizers,verbs=update

func (r *SearchController) Reconcile(ctx common.Context, search *v1beta1.Search) error {

	stack, err := stacks.GetStack(ctx, search.Spec)
	if err != nil {
		return err
	}

	elasticSearchConfiguration, err := elasticsearchconfigurations.Require(ctx, stack.Name)
	if err != nil {
		return err
	}

	env := []corev1.EnvVar{
		common.Env("OPEN_SEARCH_SERVICE", fmt.Sprintf("%s:%d", elasticSearchConfiguration.Spec.Host, elasticSearchConfiguration.Spec.Port)),
		common.Env("OPEN_SEARCH_SCHEME", elasticSearchConfiguration.Spec.Scheme),
		common.Env("ES_INDICES", "stacks"),
	}
	if elasticSearchConfiguration.Spec.BasicAuth != nil {
		if elasticSearchConfiguration.Spec.BasicAuth.SecretName == "" {
			env = append(env,
				common.Env("OPEN_SEARCH_USERNAME", elasticSearchConfiguration.Spec.BasicAuth.Username),
				common.Env("OPEN_SEARCH_PASSWORD", elasticSearchConfiguration.Spec.BasicAuth.Password),
			)
		} else {
			env = append(env,
				common.EnvFromSecret("OPEN_SEARCH_USERNAME", elasticSearchConfiguration.Spec.BasicAuth.SecretName, "username"),
				common.EnvFromSecret("OPEN_SEARCH_PASSWORD", elasticSearchConfiguration.Spec.BasicAuth.SecretName, "password"),
			)
		}
	}

	image := common.GetImage("search", common.GetVersion(stack, search.Spec.Version))
	_, _, err = common.CreateOrUpdate[*v1beta1.StreamProcessor](ctx, types.NamespacedName{
		Name: common.GetObjectName(stack.Name, "stream-processor"),
	},
		common.WithController[*v1beta1.StreamProcessor](ctx.GetScheme(), search),
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

	_, _, err = common.CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "search",
	},
		common.WithController[*appsv1.Deployment](ctx.GetScheme(), search),
		deployments.WithMatchingLabels("search"),
		deployments.WithContainers(corev1.Container{
			Name:            "search",
			Image:           common.GetImage("search", common.GetVersion(stack, search.Spec.Version)),
			Ports:           []corev1.ContainerPort{deployments.StandardHTTPPort()},
			Env:             env,
			Resources:       common.GetResourcesWithDefault(search.Spec.ResourceProperties, common.ResourceSizeSmall()),
			ImagePullPolicy: common.GetPullPolicy(image),
		}),
	)

	if err := httpapis.Create(ctx, stack, search, "search"); err != nil {
		return err
	}

	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SearchController) SetupWithManager(mgr common.Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Search{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&v1beta1.StreamProcessor{}).
		Owns(&v1beta1.HTTPAPI{}).
		Owns(&appsv1.Deployment{}), nil
}

func ForSearch() *SearchController {
	return &SearchController{}
}
