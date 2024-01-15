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
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/stacks"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// AuthClientController reconciles a Auth object
type AuthClientController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=authclients,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=authclients/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=authclients/finalizers,verbs=update

func (r *AuthClientController) Reconcile(ctx core.Context, authClient *v1beta1.AuthClient) error {

	stack, err := stacks.GetStack(ctx, authClient)
	if err != nil {
		return err
	}

	_, _, err = core.CreateOrUpdate[*corev1.Secret](ctx, types.NamespacedName{
		Name:      fmt.Sprintf("auth-client-%s", authClient.Name),
		Namespace: stack.Name,
	},
		func(t *corev1.Secret) {
			t.StringData = map[string]string{
				"id":     authClient.Spec.ID,
				"secret": authClient.Spec.Secret,
			}
		},
		core.WithController[*corev1.Secret](ctx.GetScheme(), authClient),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *AuthClientController) SetupWithManager(mgr core.Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.AuthClient{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Secret{}), nil
}

func ForAuthClient() *AuthClientController {
	return &AuthClientController{}
}
