package stacks

import (
	"context"
	"github.com/formancehq/operator/v2/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func WatchDependents(mgr core.Manager, list client.ObjectList) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {

		stack, err := StackForDependent(object)
		if err != nil {
			return nil
		}
		objects, err := GetDependentObjects(core.NewContext(mgr, ctx), stack, list)
		if err != nil {
			return nil
		}

		return core.MapObjectToReconcileRequests(objects...)
	}
}
