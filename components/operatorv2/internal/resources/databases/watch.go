package databases

import (
	"context"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch(mgr core.Manager, service string, list client.ObjectList) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		database := object.(*v1beta1.Database)
		if database.Spec.Service != service {
			return []reconcile.Request{}
		}

		list = list.DeepCopyObject().(client.ObjectList)
		list, err := stacks.GetDependentObjects(core.NewContext(mgr, ctx), database.Spec.Stack, list)
		if err != nil {
			return []reconcile.Request{}
		}

		return core.MapObjectToReconcileRequests(list...)
	}
}
