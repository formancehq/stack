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
	"context"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	v1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SearchController reconciles a Search object
type SearchController struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=searches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=searches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=searches/finalizers,verbs=update

func (r *SearchController) Reconcile(ctx context.Context, search *v1beta1.Search) error {

	//stack, err := GetStack(ctx, r.Client, search.Spec)
	//if err != nil {
	//	return err
	//}
	//
	//elasticSearchConfiguration, err := RequireElasticSearchConfiguration(ctx, r.Client, stack.Name)
	//if err != nil {
	//	return err
	//}
	//
	//CreateOrUpdate[*v1beta1.StreamProcessor](ctx, r.Client, types.NamespacedName{
	//	Name:      GetObjectName(stack.Name, "search"),
	//},
	//	WithController[*v1beta1.StreamProcessor](r.Scheme, search),
	//	func(t *v1beta1.StreamProcessor) {
	//		t.Spec.Stack = stack.Name
	//		t.Spec.Batching = search.Spec.Batching
	//		t.Spec.ResourceProperties = search.Spec.ResourceProperties
	//		t.Spec.DevProperties = search.Spec.DevProperties
	//	},
	//)

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SearchController) SetupWithManager(mgr ctrl.Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Search{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForSearch(client client.Client, scheme *runtime.Scheme) *SearchController {
	return &SearchController{
		Client: client,
		Scheme: scheme,
	}
}
