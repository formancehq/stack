package topics

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
		topic := object.(*v1beta1.Topic)
		if topic.Spec.Service != service {
			return []reconcile.Request{}
		}

		objects, ret := stacks.GetDependentObjects(core.NewContext(mgr, ctx), topic.Spec.Stack, list)
		if ret != nil {
			return []reconcile.Request{}
		}

		return core.MapObjectToReconcileRequests(objects...)
	}
}
