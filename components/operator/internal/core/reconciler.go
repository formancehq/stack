package core

import (
	"context"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	pkgError "github.com/pkg/errors"
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
	return collectionutils.Map(items, func(object T) reconcile.Request {
		return reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      object.GetName(),
				Namespace: object.GetNamespace(),
			},
		}
	})
}

func Setup(mgr ctrl.Manager, platform Platform) error {
	for _, rtype := range mgr.GetScheme().AllKnownTypes() {

		object, ok := reflect.New(rtype).Interface().(client.Object)
		if !ok {
			continue
		}

		_, ok = object.(Dependent)
		if !ok {
			continue
		}

		mgr.GetLogger().Info("Detect stack dependency object, automatically index field", "type", rtype)
		if err := mgr.GetFieldIndexer().
			IndexField(context.Background(), object, "stack", func(object client.Object) []string {
				return []string{object.(Dependent).GetStack()}
			}); err != nil {
			mgr.GetLogger().Error(err, "indexing stack field", "type", rtype)
			return err
		}

		kinds, _, err := mgr.GetScheme().ObjectKinds(object)
		if err != nil {
			return err
		}
		us := &unstructured.Unstructured{}
		us.SetGroupVersionKind(kinds[0])
		if err := mgr.GetFieldIndexer().
			IndexField(context.Background(), us, "stack", func(object client.Object) []string {
				stack := object.(*unstructured.Unstructured).Object["spec"].(map[string]any)["stack"]
				if stack == nil {
					return []string{}
				}
				return []string{stack.(string)}
			}); err != nil {
			mgr.GetLogger().Error(err, "indexing stack field", "type", &unstructured.Unstructured{})
			return err
		}
	}

	wrappedMgr := NewDefaultManager(mgr, platform)
	for _, initializer := range initializers {
		if err := initializer(wrappedMgr); err != nil {
			return err
		}
	}

	return nil
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

func WithWatchConfigurationObject(t client.Object) reconcilerOption {
	return func(mgr Manager, builder *builder.Builder, target client.Object) error {
		builder.Watches(t, handler.EnqueueRequestsFromMapFunc(WatchUsingLabels(mgr, target)))

		return nil
	}
}

func WithWatchDependency(t Dependent) reconcilerOption {
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
			return fn(NewContext(mgr, ctx), object.(T))
		}))

		return nil
	}
}

func WithReconciler[T Object](controller Controller[T], opts ...reconcilerOption) Initializer {
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

			setStatus := func(err error) {
				if err != nil {
					object.SetReady(false)
					object.SetError(err.Error())
				} else {
					object.SetReady(true)
					object.SetError("")
				}
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
				setStatus(err)
				if !pkgError.Is(err, ErrPending) &&
					!pkgError.Is(err, ErrDeleted) {
					reconcilerError = err
				}
			} else {
				setStatus(nil)
			}

			if err := mgr.GetClient().Status().Patch(ctx, object, patch); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, reconcilerError
		}))
	}
}

func WithStackDependencyReconciler[T Dependent](fn func(ctx Context, t T) error, opts ...reconcilerOption) Initializer {
	return WithReconciler(ForStackDependency(fn), opts...)
}

func WithModuleReconciler[T Module](fn func(ctx Context, t T) error, opts ...reconcilerOption) Initializer {
	return WithStackDependencyReconciler(ForModule(fn), opts...)
}

func WithIndex[T client.Object](name string, eval func(t T) string) Initializer {
	return func(mgr Manager) error {
		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		return mgr.GetFieldIndexer().
			IndexField(context.Background(), t, name, func(rawObj client.Object) []string {
				return []string{
					eval(rawObj.(T)),
				}
			})
	}
}
