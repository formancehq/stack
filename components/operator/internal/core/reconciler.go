package core

import (
	"context"
	"reflect"
	"strings"

	"k8s.io/client-go/util/workqueue"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func MapObjectToReconcileRequests[T client.Object](items ...T) []reconcile.Request {
	return Map(items, func(object T) reconcile.Request {
		return reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      object.GetName(),
				Namespace: object.GetNamespace(),
			},
		}
	})
}

type Initializer func(mgr Manager) error

var initializers = make([]Initializer, 0)

func Init(i ...Initializer) {
	initializers = append(initializers, i...)
}

type reconcilerOption func(mgr Manager, builder *builder.Builder, target client.Object) error

func WithOwn(v client.Object, opts ...builder.OwnsOption) reconcilerOption {
	return func(mgr Manager, builder *builder.Builder, target client.Object) error {
		builder.Owns(v, opts...)

		return nil
	}
}

func WithWatchSettings() reconcilerOption {

	buildReconcileRequests := func(ctx context.Context, mgr Manager, target client.Object, opts ...client.ListOption) []reconcile.Request {
		kinds, _, err := mgr.GetScheme().ObjectKinds(target)
		if err != nil {
			return []reconcile.Request{}
		}

		us := &unstructured.UnstructuredList{}
		us.SetGroupVersionKind(kinds[0])
		if err := mgr.GetClient().List(ctx, us, opts...); err != nil {
			return []reconcile.Request{}
		}

		return MapObjectToReconcileRequests(
			Map(us.Items, ToPointer[unstructured.Unstructured])...,
		)
	}

	return func(mgr Manager, builder *builder.Builder, target client.Object) error {
		builder.Watches(&v1beta1.Settings{}, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
			settings := object.(*v1beta1.Settings)

			ret := make([]reconcile.Request, 0)
			if !settings.IsWildcard() {
				for _, stack := range settings.GetStacks() {
					ret = append(ret, buildReconcileRequests(ctx, mgr, target, client.MatchingFields{
						"stack": stack,
					})...)
				}
			} else {
				ret = append(ret, buildReconcileRequests(ctx, mgr, target)...)
			}

			return ret
		}))

		return nil
	}
}

func WithWatchDependency(t v1beta1.Dependent) reconcilerOption {
	return func(mgr Manager, builder *builder.Builder, target client.Object) error {
		builder.Watches(t, handler.EnqueueRequestsFromMapFunc(WatchDependents(mgr, target)))

		return nil
	}
}

func WithWatchStack() reconcilerOption {
	return func(mgr Manager, builder *builder.Builder, target client.Object) error {
		builder.Watches(&v1beta1.Stack{}, handler.EnqueueRequestsFromMapFunc(Watch(mgr, target)))

		return nil
	}
}

func WithWatch[T client.Object](fn func(ctx Context, object T) []reconcile.Request) reconcilerOption {
	return func(mgr Manager, builder *builder.Builder, target client.Object) error {
		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		builder.Watches(t, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
			return fn(NewContext(mgr.GetClient(), mgr.GetScheme(), mgr.GetPlatform(), ctx), object.(T))
		}))

		return nil
	}
}

func WithReconciler[T client.Object](controller Controller[T], opts ...reconcilerOption) Initializer {
	return func(mgr Manager) error {

		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		b := ctrl.NewControllerManagedBy(mgr).
			For(t, builder.WithPredicates(predicate.Or(
				predicate.GenerationChangedPredicate{},
				predicate.Funcs{
					CreateFunc: func(event event.CreateEvent) bool {
						return true
					},
					DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
						return true
					},
					UpdateFunc: func(updateEvent event.UpdateEvent) bool {
					l:
						for _, referenceFromNew := range updateEvent.ObjectNew.GetOwnerReferences() {
							for _, referenceFromOld := range updateEvent.ObjectOld.GetOwnerReferences() {
								if referenceFromNew.UID == referenceFromOld.UID {
									continue l
								}
							}
							return true
						}

						return len(updateEvent.ObjectOld.GetOwnerReferences()) != len(updateEvent.ObjectNew.GetOwnerReferences())
					},
					GenericFunc: func(genericEvent event.GenericEvent) bool {
						return true
					},
				},
			)))
		for _, opt := range opts {
			if err := opt(mgr, b, t); err != nil {
				return err
			}
		}

		return b.Complete(reconcile.Func(func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

			var object T
			object = reflect.New(reflect.TypeOf(object).Elem()).Interface().(T)
			if err := mgr.GetClient().Get(ctx, types.NamespacedName{
				Name: request.Name,
			}, object); err != nil {
				if errors.IsNotFound(err) {
					return ctrl.Result{}, nil
				}
				return ctrl.Result{}, err
			}

			cp := object.DeepCopyObject().(T)
			patch := client.MergeFrom(cp)

			var reconcilerError error
			err := controller(struct {
				context.Context
				Manager
			}{
				Context: ctx,
				Manager: mgr,
			}, object)
			if err != nil {
				reconcilerError = err
			}

			if err := mgr.GetClient().Status().Patch(ctx, object, patch); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, reconcilerError
		}))
	}
}

func WithStdReconciler[T v1beta1.Object](ctrl func(ctx Context, t T) error, opts ...reconcilerOption) Initializer {
	return WithReconciler(ForReadier(ctrl), opts...)
}

func WithStackDependencyReconciler[T v1beta1.Dependent](fn func(ctx Context, stack *v1beta1.Stack, t T) error, opts ...reconcilerOption) Initializer {
	opts = append(opts, WithWatchStack())
	return WithStdReconciler(ForStackDependency(fn), opts...)
}

func WithModuleReconciler[T v1beta1.Module](fn func(ctx Context, stack *v1beta1.Stack, t T, version string) error, opts ...reconcilerOption) Initializer {
	opts = append(opts, WithWatchVersions)
	return WithStackDependencyReconciler(ForModule(fn), opts...)
}

func WithWatchVersions(mgr Manager, builder *builder.Builder, target client.Object) error {
	reconcileModule := func(ctx context.Context, versionFileName string, limitingInterface workqueue.RateLimitingInterface) {
		stackList := &v1beta1.StackList{}
		if err := mgr.GetClient().List(ctx, stackList, client.MatchingFields{
			".spec.versionsFromFile": versionFileName,
		}); err != nil {
			panic(err)
		}

		kinds, _, err := mgr.GetScheme().ObjectKinds(target)
		if err != nil {
			panic(err)
		}

		for _, stack := range stackList.Items {
			list := &unstructured.UnstructuredList{}
			list.SetGroupVersionKind(kinds[0])
			if err := mgr.GetClient().List(ctx, list, client.MatchingFields{
				"stack": stack.Name,
			}); err != nil {
				panic(err)
			}

			for _, item := range list.Items {
				limitingInterface.Add(reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name: item.GetName(),
					},
				})
			}
		}
	}
	builder.Watches(&v1beta1.Versions{}, handler.Funcs{
		CreateFunc: func(ctx context.Context, createEvent event.CreateEvent, limitingInterface workqueue.RateLimitingInterface) {
			reconcileModule(ctx, createEvent.Object.GetName(), limitingInterface)
		},
		UpdateFunc: func(ctx context.Context, updateEvent event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
			oldObject := updateEvent.ObjectOld.(*v1beta1.Versions)
			newObject := updateEvent.ObjectNew.(*v1beta1.Versions)

			kinds, _, err := mgr.GetScheme().ObjectKinds(target)
			if err != nil {
				panic(err)
			}
			kind := strings.ToLower(kinds[0].Kind)
			if oldObject.Spec[kind] == newObject.Spec[kind] {
				return
			}

			reconcileModule(ctx, updateEvent.ObjectNew.GetName(), limitingInterface)
		},
		DeleteFunc: func(ctx context.Context, deleteEvent event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
			reconcileModule(ctx, deleteEvent.Object.GetName(), limitingInterface)
		},
	})

	return nil
}

func WithIndex[T client.Object](name string, eval func(t T) []string) Initializer {
	return func(mgr Manager) error {
		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		return mgr.GetFieldIndexer().
			IndexField(context.Background(), t, name, func(rawObj client.Object) []string {
				return eval(rawObj.(T))
			})
	}
}

func WithSimpleIndex[T client.Object](name string, eval func(t T) string) Initializer {
	return WithIndex(name, func(t T) []string {
		return []string{eval(t)}
	})
}
