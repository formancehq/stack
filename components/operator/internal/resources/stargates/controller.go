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

package stargates

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	appsv1 "k8s.io/api/apps/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=stargates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=stargates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=stargates/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, stargate *v1beta1.Stargate, version string) error {
	if err := createDeployment(ctx, stack, stargate, version); err != nil {
		return err
	}

	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithWatchSettings[*v1beta1.Stargate](),
			WithOwn[*v1beta1.Stargate](&appsv1.Deployment{}),
		),
	)
}
