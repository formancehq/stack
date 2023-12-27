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

package controller

import (
	"context"

	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
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
	gcTopicsFinalizer = "gc-topics"
)

// TopicQueryReconciler reconciles a TopicQuery object
type TopicQueryReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=topicqueries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=topicqueries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=topicqueries/finalizers,verbs=update

func (r *TopicQueryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx, "topicquery", req.NamespacedName)
	log.Info("Starting reconciliation")

	topicQuery := &v1beta1.TopicQuery{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: req.Name,
	}, topicQuery); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !topicQuery.DeletionTimestamp.IsZero() {
		topic := &v1beta1.Topic{}
		if err := r.Client.Get(ctx, types.NamespacedName{
			Name: GetObjectName(topicQuery.Spec.Stack, topicQuery.Spec.Service),
		}, topic); err != nil {
			if !errors.IsNotFound(err) {
				return ctrl.Result{}, err
			}
		} else {
			topic.Spec.Queries = Filter(topic.Spec.Queries, func(s string) bool {
				return s != topicQuery.Spec.QueriedBy
			})
			if err := r.Client.Update(ctx, topic); err != nil {
				return ctrl.Result{}, err
			}
		}

		if updated := controllerutil.RemoveFinalizer(topicQuery, gcTopicsFinalizer); updated {
			if err := r.Client.Update(ctx, topicQuery); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if updated := controllerutil.AddFinalizer(topicQuery, gcTopicsFinalizer); updated {
		if err := r.Client.Update(ctx, topicQuery); err != nil {
			return ctrl.Result{}, err
		}
	}

	topic := &v1beta1.Topic{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: GetObjectName(topicQuery.Spec.Stack, topicQuery.Spec.Service),
	}, topic); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		if err := r.Client.Create(ctx, &v1beta1.Topic{
			ObjectMeta: ctrl.ObjectMeta{
				Name: GetObjectName(topicQuery.Spec.Stack, topicQuery.Spec.Service),
			},
			Spec: v1beta1.TopicSpec{
				Queries: []string{topicQuery.Spec.QueriedBy},
				Stack:   topicQuery.Spec.Stack,
				Service: topicQuery.Spec.Service,
			},
		}); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if !Contains(topic.Spec.Queries, topicQuery.Spec.QueriedBy) {
		topic.Spec.Queries = append(topic.Spec.Queries, topicQuery.Spec.QueriedBy)
		if err := r.Client.Update(ctx, topic); err != nil {
			return ctrl.Result{}, err
		}
	}

	if topicQuery.Status.Ready != topic.Status.Ready {
		topicQuery.Status.Ready = topic.Status.Ready
		if err := r.Client.Status().Update(ctx, topicQuery); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TopicQueryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.TopicQuery{}, ".spec.service", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.TopicQuery).Spec.Service}
	}); err != nil {
		return err
	}

	if err := indexer.IndexField(context.Background(), &v1beta1.TopicQuery{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.TopicQuery).Spec.Stack}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.TopicQuery{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1beta1.Topic{},
			// Watch update of Topic to be able to set Ready flag on TopicQuery
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				list := v1beta1.TopicQueryList{}
				if err := mgr.GetClient().List(ctx, &list, client.MatchingFields{
					".spec.service": object.(*v1beta1.Topic).Spec.Service,
					".spec.stack":   object.(*v1beta1.Topic).Spec.Stack,
				}); err != nil {
					log.FromContext(ctx).Error(err, "listing topic queries")
					return []reconcile.Request{}
				}

				return MapObjectToReconcileRequests(
					Map(list.Items, ToPointer[v1beta1.TopicQuery])...,
				)
			}),
		).
		Complete(r)
}

func NewTopicQueryReconciler(client client.Client, scheme *runtime.Scheme) *TopicQueryReconciler {
	return &TopicQueryReconciler{
		Client: client,
		Scheme: scheme,
	}
}
