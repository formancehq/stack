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
	"github.com/formancehq/operator/internal/resources/services"
	corev1 "k8s.io/api/core/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=httpapis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=httpapis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=httpapis/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, httpAPI *v1beta1.HTTPAPI) error {
	_, err := services.Create(ctx, httpAPI, httpAPI.Spec.Name)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	Init(
		WithStackDependencyReconciler(Reconcile,
			WithOwn(&corev1.Service{}),
			WithWatchConfigurationObject(&v1beta1.Settings{}),
		),
	)
}
