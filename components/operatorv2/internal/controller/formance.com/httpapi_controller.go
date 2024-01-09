/*
Copyright 2022.

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
	v1beta1 "github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// HTTPAPI reconciles a HTTPAPI object
type HTTPAPI struct{}

//+kubebuilder:rbac:groups=formance.com,resources=httpapis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=httpapis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=httpapis/finalizers,verbs=update

func (r *HTTPAPI) Reconcile(ctx Context, httpAPI *v1beta1.HTTPAPI) error {
	_, operationResult, err := CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Namespace: httpAPI.Spec.Stack,
		Name:      httpAPI.Spec.Name,
	},
		func(t *corev1.Service) {
			t.ObjectMeta.Annotations = httpAPI.Spec.Annotations

			t.Labels = map[string]string{
				"app.kubernetes.io/service-name": httpAPI.Name,
			}
			t.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:        "http",
					Port:        8080,
					Protocol:    "TCP",
					AppProtocol: pointer.String("http"),
					TargetPort:  intstr.FromString("http"),
				}},
				Selector: map[string]string{
					"app.kubernetes.io/name": httpAPI.Spec.Name,
				},
			}
		},
		WithController[*corev1.Service](ctx.GetScheme(), httpAPI),
	)
	if err != nil {
		httpAPI.Status.SetCondition(v1beta1.Condition{
			Type:               "ServiceExists",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: httpAPI.Generation,
			LastTransitionTime: metav1.Now(),
			Message:            err.Error(),
		})
		httpAPI.Status.Ready = false
		return nil
	}
	httpAPI.Status.Ready = true

	switch operationResult {
	case controllerutil.OperationResultNone:
	case controllerutil.OperationResultCreated:
		httpAPI.Status.SetCondition(v1beta1.Condition{
			Type:               "ServiceExists",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: httpAPI.Generation,
			LastTransitionTime: metav1.Now(),
		})
	case controllerutil.OperationResultUpdated:
		httpAPI.Status.SetCondition(v1beta1.Condition{
			Type:               "ServiceExists",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: httpAPI.Generation,
			LastTransitionTime: metav1.Now(),
		})
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPAPI) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.HTTPAPI{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Service{}), nil
}

func ForHTTPAPI() *HTTPAPI {
	return &HTTPAPI{}
}
