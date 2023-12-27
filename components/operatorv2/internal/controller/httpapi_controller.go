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

package controller

import (
	"context"

	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HTTPAPIReconciler reconciles a HTTPAPI object
type HTTPAPIReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=httpapis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=httpapis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=httpapis/finalizers,verbs=update

func (r *HTTPAPIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx, "HTTPAPI", req.NamespacedName)
	log.Info("Starting reconciliation")

	HTTPAPI := &v1beta1.HTTPAPI{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: req.Name,
	}, HTTPAPI); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	_, _, err := CreateOrUpdate[*corev1.Service](ctx, r.Client, types.NamespacedName{
		Namespace: HTTPAPI.Spec.Stack,
		Name:      HTTPAPI.Spec.Name,
	},
		func(t *corev1.Service) {
			t.ObjectMeta.Annotations = HTTPAPI.Spec.Annotations

			t.Labels = map[string]string{
				"app.kubernetes.io/service-name": HTTPAPI.Name,
			}
			t.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:        "http",
					Port:        8080,
					Protocol:    "TCP",
					AppProtocol: pointer.String("http"),
					TargetPort:  intstr.FromString(HTTPAPI.Spec.PortName),
				}},
				Selector: map[string]string{
					"app.kubernetes.io/name": HTTPAPI.Spec.Name,
				},
			}
		},
		WithController[*corev1.Service](r.Scheme, HTTPAPI),
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPAPIReconciler) SetupWithManager(mgr ctrl.Manager) error {

	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.HTTPAPI{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.HTTPAPI).Spec.Stack}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.HTTPAPI{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func NewHTTPAPIReconciler(client client.Client, scheme *runtime.Scheme) *HTTPAPIReconciler {
	return &HTTPAPIReconciler{
		Client: client,
		Scheme: scheme,
	}
}
