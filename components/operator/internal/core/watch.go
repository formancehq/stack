package core

import (
	"context"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func WatchDependents[T client.Object](mgr Manager) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {

		objects, err := GetAllDependents[T](NewContext(mgr, ctx), object.(Dependent).GetStack())
		if err != nil {
			return nil
		}

		return MapObjectToReconcileRequests(objects...)
	}
}

func Watch[T client.Object](mgr Manager) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {

		objects, err := GetAllDependents[T](NewContext(mgr, ctx), object.(*v1beta1.Stack).Name)
		if err != nil {
			return nil
		}

		return MapObjectToReconcileRequests(objects...)
	}
}
