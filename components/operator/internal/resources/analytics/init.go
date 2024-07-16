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

package analytics

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
)

//+kubebuilder:rbac:groups=formance.com,resources=analytics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=analytics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=analytics/finalizers,verbs=update

func Reconcile(_ Context, _ *v1beta1.Stack, _ *v1beta1.Analytics, _ string) error {
	return nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile),
	)
}
