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

package wallets

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/authclients"
	"github.com/formancehq/operator/internal/resources/gatewayhttpapis"
	appsv1 "k8s.io/api/apps/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=wallets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=wallets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=wallets/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, wallets *v1beta1.Wallets, version string) error {

	hasAuth, err := HasDependency(ctx, wallets.Spec.Stack, &v1beta1.Auth{})
	if err != nil {
		return err
	}
	var authClient *v1beta1.AuthClient
	if hasAuth {
		authClient, err = authclients.Create(ctx, stack, wallets, "wallets", func(spec *v1beta1.AuthClientSpec) {
			spec.Scopes = []string{"ledger:read", "ledger:write"}
		})
		if err != nil {
			return err
		}
	}

	if err := createDeployment(ctx, stack, wallets, authClient, version); err != nil {
		return err
	}

	if err := gatewayhttpapis.Create(ctx, wallets); err != nil {
		return err
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithWatchSettings[*v1beta1.Wallets](),
			WithWatchDependency[*v1beta1.Wallets](&v1beta1.Auth{}),
			WithOwn[*v1beta1.Wallets](&v1beta1.AuthClient{}),
			WithOwn[*v1beta1.Wallets](&appsv1.Deployment{}),
			WithOwn[*v1beta1.Wallets](&v1beta1.GatewayHTTPAPI{}),
		),
	)
}
