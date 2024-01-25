package core

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func buildReconcileRequests(ctx context.Context, mgr Manager, target client.Object, opts ...client.ListOption) []reconcile.Request {
	kinds, _, err := mgr.GetScheme().ObjectKinds(target)
	if err != nil {
		return []reconcile.Request{}
	}

	us := &unstructured.UnstructuredList{}
	us.SetGroupVersionKind(kinds[0])
	if err := mgr.GetClient().List(ctx, us, opts...); err != nil {
		return []reconcile.Request{}
	}

	return MapObjectToReconcileRequests(
		collectionutils.Map(us.Items, collectionutils.ToPointer[unstructured.Unstructured])...,
	)
}

func WatchConfigurationObject(mgr Manager, target client.Object) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {

		configurationObject := object.(v1beta1.ConfigurationObject)

		ret := make([]reconcile.Request, 0)
		if !configurationObject.IsWildcard() {
			for _, stack := range configurationObject.GetStacks() {
				ret = append(ret, buildReconcileRequests(ctx, mgr, target, client.MatchingFields{
					"stack": stack,
				})...)
			}
		} else {
			ret = append(ret, buildReconcileRequests(ctx, mgr, target)...)
		}

		return ret
	}
}
