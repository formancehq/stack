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

package brokertopics

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokers"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
)

//+kubebuilder:rbac:groups=formance.com,resources=brokertopics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=brokertopics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=brokertopics/finalizers,verbs=update

func Reconcile(ctx core.Context, stack *v1beta1.Stack, topic *v1beta1.BrokerTopic) error {

	if len(topic.GetOwnerReferences()) == 1 { // Remains only the stack
		return ctx.GetClient().Delete(ctx, topic)
	}

	broker, _, err := core.CreateOrUpdate[*v1beta1.Broker](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.Broker) error {
		t.Spec.Stack = stack.Name
		return nil
	}, core.WithController[*v1beta1.Broker](ctx.GetScheme(), stack))
	if err != nil {
		return err
	}

	if !broker.Status.Ready {
		return core.NewApplicationError().WithMessage("broker not ready")
	}
	topic.Status.Ready = broker.Status.Ready

	return nil
}

func init() {
	core.Init(
		core.WithResourceReconciler(Reconcile,
			core.WithOwn[*v1beta1.BrokerTopic](&batchv1.Job{}),
			core.WithWatchSettings[*v1beta1.BrokerTopic](),
			brokers.Watch[*v1beta1.BrokerTopic](),
		),
		core.WithSimpleIndex[*v1beta1.BrokerTopic](".spec.service", func(t *v1beta1.BrokerTopic) string {
			return t.Spec.Service
		}),
	)
}
