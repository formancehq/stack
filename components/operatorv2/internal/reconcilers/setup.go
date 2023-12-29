package reconcilers

import (
	"context"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/stacks"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type reconciler interface {
	SetupWithManager(mgr core.Manager) error
}

var reconcilers []reconciler

func Register(newReconcilers ...reconciler) {
	reconcilers = append(reconcilers, newReconcilers...)
}

func Setup(mgr ctrl.Manager, platform core.Platform) error {
	wrappedMgr := core.NewDefaultManager(mgr, platform)
	for _, reconciler := range reconcilers {
		if err := reconciler.SetupWithManager(wrappedMgr); err != nil {
			return err
		}
	}
	for _, rtype := range mgr.GetScheme().AllKnownTypes() {

		object, ok := reflect.New(rtype).Interface().(client.Object)
		if !ok {
			continue
		}

		isStackDependent, err := stacks.IsStackDependent(object)
		if err != nil {
			return err
		}

		if isStackDependent {
			mgr.GetLogger().Info("Detect stack dependency object, automatically index field", "type", rtype)
			if err := mgr.GetFieldIndexer().
				IndexField(context.Background(), reflect.New(rtype).Interface().(client.Object), ".spec.stack", func(object client.Object) []string {
					stackForDependent, err := stacks.StackForDependent(object)
					if err != nil {
						return nil
					}
					return []string{stackForDependent}
				}); err != nil {
				mgr.GetLogger().Error(err, "indexing .spec.stack field", "type", rtype)
				return err
			}
		}
	}
	return nil
}
