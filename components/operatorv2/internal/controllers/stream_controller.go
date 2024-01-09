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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// StreamController reconciles a Stream object
type StreamController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=streams,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=streams/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=streams/finalizers,verbs=update

func (r *StreamController) Reconcile(ctx Context, stream *v1beta1.Stream) error {
	_, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stream.Spec.Stack,
		Name:      fmt.Sprintf("stream-%s", stream.Name),
	},
		WithController[*corev1.ConfigMap](ctx.GetScheme(), stream),
		func(t *corev1.ConfigMap) {
			t.Data = map[string]string{
				"stream.yaml": stream.Spec.Data,
			}
		},
	)
	if err != nil {
		return err
	}

	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *StreamController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		Owns(&corev1.ConfigMap{}).
		For(&v1beta1.Stream{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForStream() *StreamController {
	return &StreamController{}
}
