package databases

import (
	"context"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch[LIST client.ObjectList, OBJECT client.Object](mgr core.Manager, service string) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		database := object.(*v1beta1.Database)
		if database.Spec.Service != service {
			return []reconcile.Request{}
		}

		list, err := stacks.GetStackDependencyObjects[LIST, OBJECT](core.NewDefaultContext(mgr, ctx), database.Spec.Stack)
		if err != nil {
			return []reconcile.Request{}
		}

		return core.MapObjectToReconcileRequests(list...)
	}
}
