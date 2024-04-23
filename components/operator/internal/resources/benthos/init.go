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

package benthos

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	Init(
		WithStackDependencyReconciler(Reconcile,
			WithWatchSettings[*v1beta1.Benthos](),
			WithWatchDependency[*v1beta1.Benthos](&v1beta1.BenthosStream{}),
			WithOwn[*v1beta1.Benthos](&corev1.ConfigMap{}),
			WithOwn[*v1beta1.Benthos](&appsv1.Deployment{}),
			WithOwn[*v1beta1.Benthos](&v1beta1.ResourceReference{}),
			brokers.Watch[*v1beta1.Benthos](),
		),
	)
}
