package reconcilers

import (
	"context"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Context interface {
	context.Context
	Manager
}

type Controller[T client.Object] interface {
	Reconcile(ctx Context, req T) error
	SetupWithManager(mgr Manager) (*ctrl.Builder, error)
}

type reconciler interface {
	SetupWithManager(mgr Manager) error
}

var reconcilers []reconciler

func Register(newReconcilers ...reconciler) {
	reconcilers = append(reconcilers, newReconcilers...)
}

type stackDependent interface {
	GetStack() string
}

func Setup(mgr ctrl.Manager, platform Platform) error {
	wrappedMgr := newDefaultManager(mgr, platform)
	for _, reconciler := range reconcilers {
		if err := reconciler.SetupWithManager(wrappedMgr); err != nil {
			return err
		}
	}
	for _, rtype := range mgr.GetScheme().AllKnownTypes() {
		specField, ok := rtype.FieldByName("Spec")
		if !ok {
			continue
		}

		if specField.Type.AssignableTo(reflect.TypeOf((*stackDependent)(nil)).Elem()) {
			mgr.GetLogger().Info("Detect stack dependency object, automatically index field", "type", rtype)
			if err := mgr.GetFieldIndexer().
				IndexField(context.Background(), reflect.New(rtype).Interface().(client.Object), ".spec.stack", func(object client.Object) []string {
					return []string{
						reflect.ValueOf(object).
							Elem().
							FieldByName("Spec").
							Interface().(stackDependent).
							GetStack(),
					}
				}); err != nil {
				mgr.GetLogger().Error(err, "indexing .spec.stack field", "type", rtype)
				return err
			}
		}
	}
	return nil
}
