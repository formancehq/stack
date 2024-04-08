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

package brokerconsumers

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//+kubebuilder:rbac:groups=formance.com,resources=brokerconsumers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=brokerconsumers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=brokerconsumers/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, topicQuery *v1beta1.BrokerConsumer) error {

	for _, service := range topicQuery.Spec.Services {
		topic := &v1beta1.BrokerTopic{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: GetObjectName(topicQuery.Spec.Stack, service),
		}, topic); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			topic = &v1beta1.BrokerTopic{
				ObjectMeta: ctrl.ObjectMeta{
					Name: GetObjectName(topicQuery.Spec.Stack, service),
				},
				Spec: v1beta1.BrokerTopicSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: topicQuery.Spec.Stack,
					},
					Service: service,
				},
			}
			if err := controllerutil.SetOwnerReference(topicQuery, topic, ctx.GetScheme()); err != nil {
				return err
			}

			if err := controllerutil.SetOwnerReference(stack, topic, ctx.GetScheme()); err != nil {
				return err
			}
			if err := ctx.GetClient().Create(ctx, topic); err != nil {
				return err
			}
			return nil
		} else {
			patch := client.MergeFrom(topic.DeepCopy())
			if err := controllerutil.SetOwnerReference(topicQuery, topic, ctx.GetScheme()); err != nil {
				return err
			}
			if err := ctx.GetClient().Patch(ctx, topic, patch); err != nil {
				return err
			}
		}

		if !topic.Status.Ready {
			return NewPendingError()
		}
	}

	return nil
}

func init() {
	Init(
		WithStackDependencyReconciler(Reconcile,
			WithWatch[*v1beta1.BrokerConsumer, *v1beta1.BrokerTopic](func(ctx Context, object *v1beta1.BrokerTopic) []reconcile.Request {
				list := v1beta1.BrokerTopicConsumerList{}
				if err := ctx.GetClient().List(ctx, &list, client.MatchingFields{
					"queriedBy": object.Spec.Service,
					"stack":     object.Spec.Stack,
				}); err != nil {
					log.FromContext(ctx).Error(err, "listing topic queries")
					return []reconcile.Request{}
				}

				return MapObjectToReconcileRequests(
					Map(list.Items, ToPointer[v1beta1.BrokerTopicConsumer])...,
				)
			}),
		),
	)
}
