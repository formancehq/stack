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
	"context"
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	gcFinalizer = "gc"
)

// TopicQueryController reconciles a TopicQuery object
type TopicQueryController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=topicqueries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=topicqueries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=topicqueries/finalizers,verbs=update

func (r *TopicQueryController) Reconcile(ctx Context, topicQuery *v1beta1.TopicQuery) error {

	if !topicQuery.DeletionTimestamp.IsZero() {
		topic := &v1beta1.BrokerTopic{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: GetObjectName(topicQuery.Spec.Stack, topicQuery.Spec.Service),
		}, topic); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		} else {
			if err := controllerutil.RemoveOwnerReference(topicQuery, topic, ctx.GetScheme()); err != nil {
				return err
			}
			if err := ctx.GetClient().Update(ctx, topic); err != nil {
				return err
			}
		}

		if updated := controllerutil.RemoveFinalizer(topicQuery, gcFinalizer); updated {
			if err := ctx.GetClient().Update(ctx, topicQuery); err != nil {
				return err
			}
		}
		return nil
	}

	if updated := controllerutil.AddFinalizer(topicQuery, gcFinalizer); updated {
		if err := ctx.GetClient().Update(ctx, topicQuery); err != nil {
			return err
		}
	}

	topic := &v1beta1.BrokerTopic{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: GetObjectName(topicQuery.Spec.Stack, topicQuery.Spec.Service),
	}, topic); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		topic = &v1beta1.BrokerTopic{
			ObjectMeta: ctrl.ObjectMeta{
				Name: GetObjectName(topicQuery.Spec.Stack, topicQuery.Spec.Service),
			},
			Spec: v1beta1.BrokerTopicSpec{
				StackDependency: v1beta1.StackDependency{
					Stack: topicQuery.Spec.Stack,
				},
				Service: topicQuery.Spec.Service,
			},
		}
		if err := controllerutil.SetOwnerReference(topicQuery, topic, ctx.GetScheme()); err != nil {
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
		return ErrPending
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TopicQueryController) SetupWithManager(mgr Manager) (*builder.Builder, error) {

	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.TopicQuery{}, ".spec.service", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.TopicQuery).Spec.Service}
	}); err != nil {
		return nil, err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.TopicQuery{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1beta1.BrokerTopic{},
			// WatchUsingLabels update of BrokerTopic to be able to set Ready flag on TopicQuery
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				list := v1beta1.TopicQueryList{}
				if err := mgr.GetClient().List(ctx, &list, client.MatchingFields{
					".spec.service": object.(*v1beta1.BrokerTopic).Spec.Service,
					"stack":         object.(*v1beta1.BrokerTopic).Spec.Stack,
				}); err != nil {
					log.FromContext(ctx).Error(err, "listing topic queries")
					return []reconcile.Request{}
				}

				return MapObjectToReconcileRequests(
					Map(list.Items, ToPointer[v1beta1.TopicQuery])...,
				)
			}),
		), nil
}

func ForTopicQuery() *TopicQueryController {
	return &TopicQueryController{}
}
