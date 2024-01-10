package reconcilers

import (
	"context"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/stacks"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

		_, ok = object.(stacks.Dependent)
		if !ok {
			continue
		}

		mgr.GetLogger().Info("Detect stack dependency object, automatically index field", "type", rtype)
		if err := mgr.GetFieldIndexer().
			IndexField(context.Background(), object, "stack", func(object client.Object) []string {
				return []string{object.(stacks.Dependent).GetStack()}
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
	return nil
}
