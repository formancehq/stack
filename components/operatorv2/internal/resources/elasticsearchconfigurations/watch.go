package elasticsearchconfigurations

import (
	"context"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch(mgr Manager, list client.ObjectList) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		elasticSearchConfiguration := object.(*v1beta1.ElasticSearchConfiguration)

		opt := client.MatchingFields{}
		if elasticSearchConfiguration.Labels["formance.com/stack"] != "any" {
			opt[".spec.stack"] = elasticSearchConfiguration.Labels["formance.com/stack"]
		}

		if err := mgr.GetClient().List(ctx, list, opt); err != nil {
			return []reconcile.Request{}
		}

		return MapObjectToReconcileRequests(
			ExtractItemsFromList(list)...,
		)
	}
}
