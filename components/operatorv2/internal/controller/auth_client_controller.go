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
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// AuthClientReconciler reconciles a Auth object
type AuthClientReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=authclients,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=authclients/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=authclients/finalizers,verbs=update

func (r *AuthClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx, "auth-client", req.NamespacedName)
	log.Info("Starting reconciliation")

	authClient := &v1beta1.AuthClient{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: req.Name,
	}, authClient); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	stack := &v1beta1.Stack{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: authClient.Spec.Stack,
	}, stack); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	_, _, err := CreateOrUpdate[*corev1.Secret](ctx, r.Client, types.NamespacedName{
		Name:      fmt.Sprintf("auth-client-%s", authClient.Name),
		Namespace: stack.Name,
	},
		func(t *corev1.Secret) {
			t.StringData = map[string]string{
				"id":     authClient.Spec.ID,
				"secret": authClient.Spec.Secret,
			}
		},
		WithController[*corev1.Secret](r.Scheme, authClient),
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AuthClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.AuthClient{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.AuthClient).Spec.Stack}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.AuthClient{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func NewAuthClientReconciler(client client.Client, scheme *runtime.Scheme) *AuthClientReconciler {
	return &AuthClientReconciler{
		Client: client,
		Scheme: scheme,
	}
}
