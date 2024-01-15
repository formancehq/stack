package databases

import (
	"context"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch[T client.Object](mgr core.Manager, service string) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		database := object.(*v1beta1.Database)
		if database.Spec.Service != service {
			return []reconcile.Request{}
		}

		list, err := core.GetAllDependents[T](core.NewContext(mgr, ctx), database.Spec.Stack)
		if err != nil {
			return []reconcile.Request{}
		}

		return core.MapObjectToReconcileRequests(list...)
	}
}
