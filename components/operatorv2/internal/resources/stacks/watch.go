package stacks

import (
	"context"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func WatchDependency[LIST client.ObjectList, OBJECT client.Object](mgr core.Manager) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {

		stack, err := StackForDependent(object)
		if err != nil {
			return nil
		}
		objects, err := GetDependentObjects[LIST, OBJECT](core.NewContext(mgr, ctx), stack)
		if err != nil {
			return nil
		}

		return collectionutils.Map(objects, func(from OBJECT) reconcile.Request {
			return reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: from.GetNamespace(),
					Name:      from.GetName(),
				},
			}
		})
	}
}
