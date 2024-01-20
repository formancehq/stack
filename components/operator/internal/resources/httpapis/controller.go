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

package httpapis

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

//+kubebuilder:rbac:groups=formance.com,resources=httpapis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=httpapis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=httpapis/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, httpAPI *v1beta1.HTTPAPI) error {
	_, _, err := CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Namespace: httpAPI.Spec.Stack,
		Name:      httpAPI.Spec.Name,
	},
		func(t *corev1.Service) {
			if httpAPI.Spec.Service != nil {
				t.ObjectMeta.Annotations = httpAPI.Spec.Service.Annotations
			}

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
		return err
	}

	return nil
}

func init() {
	Init(
		WithStackDependencyReconciler(Reconcile,
			WithOwn(&corev1.Service{}),
		),
	)
}