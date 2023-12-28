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
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	formancev1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WalletController reconciles a Wallet object
type WalletController struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=wallets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=wallets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=wallets/finalizers,verbs=update

func (r *WalletController) Reconcile(ctx context.Context, wallet *v1beta1.Wallet) error {

	stack := &v1beta1.Stack{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: wallet.Spec.Stack,
	}, stack); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	authClient, err := CreateAuthClient(ctx, r.Client, r.Scheme, stack, wallet, "wallets", func(spec *formancev1beta1.AuthClientSpec) {
		spec.Scopes = []string{"ledger:read", "ledger:write"}
	})
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, wallet, authClient); err != nil {
		return err
	}

	if err := CreateHTTPAPI(ctx, r.Client, r.Scheme, stack, wallet, "wallets"); err != nil {
		return err
	}

	return nil
}

func (r *WalletController) createDeployment(ctx context.Context, stack *formancev1beta1.Stack, wallet *formancev1beta1.Wallet, authClient *formancev1beta1.AuthClient) error {
	env, err := GetCommonServicesEnvVars(ctx, r.Client, stack, "wallets", wallet.Spec)
	if err != nil {
		return err
	}
	env = append(env, GetAuthClientEnvVars(authClient)...)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, r.Client,
		GetNamespacedResourceName(stack.Name, "wallets"),
		func(t *appsv1.Deployment) {
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Name:      "wallets",
				Args:      []string{"serve"},
				Env:       env,
				Image:     GetImage("wallets", GetVersion(stack, wallet.Spec.Version)),
				Resources: GetResourcesWithDefault(wallet.Spec.ResourceProperties, ResourceSizeSmall()),
				Ports:     []corev1.ContainerPort{StandardHTTPPort()},
			}}
		},
		WithMatchingLabels("wallets"),
		WithController[*appsv1.Deployment](r.Scheme, wallet),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *WalletController) SetupWithManager(mgr ctrl.Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&formancev1beta1.Wallet{}), nil
}

func ForWallet(client client.Client, scheme *runtime.Scheme) *WalletController {
	return &WalletController{
		Client: client,
		Scheme: scheme,
	}
}
