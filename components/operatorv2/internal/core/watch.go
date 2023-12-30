package core

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch(mgr Manager, list client.ObjectList) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		opt := client.MatchingFields{}
		if object.GetLabels()["formance.com/stack"] != "any" {
			opt[".spec.stack"] = object.GetLabels()["formance.com/stack"]
		}

		if err := mgr.GetClient().List(ctx, list, opt); err != nil {
			return []reconcile.Request{}
		}

		return MapObjectToReconcileRequests(
			ExtractItemsFromList(list)...,
		)
	}
}
