package stacks

import (
	"context"
	"github.com/formancehq/operator/v2/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func WatchDependents[T client.Object](mgr core.Manager) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {

		objects, err := GetAllDependents[T](core.NewContext(mgr, ctx), object.(Dependent).GetStack())
		if err != nil {
			return nil
		}

		return core.MapObjectToReconcileRequests(objects...)
	}
}
