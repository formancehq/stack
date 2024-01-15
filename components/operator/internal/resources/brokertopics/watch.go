package brokertopics

import (
	"context"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch[T client.Object](mgr core.Manager, service string) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		topic := object.(*v1beta1.BrokerTopic)
		if topic.Spec.Service != service {
			return []reconcile.Request{}
		}

		objects, ret := core.GetAllDependents[T](core.NewContext(mgr, ctx), topic.Spec.Stack)
		if ret != nil {
			return []reconcile.Request{}
		}

		return core.MapObjectToReconcileRequests(objects...)
	}
}
