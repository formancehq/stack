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

package brokertopicconsumers

import (
	"fmt"
	"sort"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//+kubebuilder:rbac:groups=formance.com,resources=brokertopicconsumers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=brokertopicconsumers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=brokertopicconsumers/finalizers,verbs=update

func Reconcile(ctx Context, _ *v1beta1.Stack, topicQuery *v1beta1.BrokerTopicConsumer) error {

	l := &v1beta1.BrokerTopicConsumerList{}
	if err := ctx.GetClient().List(ctx, l, client.MatchingFields{
		"stack":     topicQuery.Spec.Stack,
		"queriedBy": topicQuery.Spec.QueriedBy,
	}); err != nil {
		return err
	}

	_, _, err := CreateOrUpdate(ctx, types.NamespacedName{
		Namespace: topicQuery.Spec.Stack,
		Name:      fmt.Sprintf("%s-%s", topicQuery.Spec.Stack, topicQuery.Spec.QueriedBy),
	}, func(t *v1beta1.BrokerConsumer) error {
		t.Spec.Stack = topicQuery.Spec.Stack
		t.Spec.QueriedBy = topicQuery.Spec.QueriedBy

	l:
		for _, expectService := range Map(l.Items, func(from v1beta1.BrokerTopicConsumer) string {
			return from.Spec.Service
		}) {
			for _, actualService := range t.Spec.Services {
				if actualService == expectService {
					continue l
				}
			}
			t.Spec.Services = append(t.Spec.Services, expectService)
		}

		sort.Strings(t.Spec.Services)

		t.SetOwnerReferences(topicQuery.OwnerReferences)

		return nil
	})
	return err
}

func init() {
	Init(
		WithStackDependencyReconciler(Reconcile),
		WithSimpleIndex[*v1beta1.BrokerTopicConsumer]("queriedBy", func(t *v1beta1.BrokerTopicConsumer) string {
			return t.Spec.QueriedBy
		}),
	)
}
