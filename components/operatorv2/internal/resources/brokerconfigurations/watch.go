package brokerconfigurations

import (
	"context"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch[LIST client.ObjectList, OBJECT client.Object](mgr Manager) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		brokerConfiguration := object.(*v1beta1.BrokerConfiguration)

		var list LIST
		list = reflect.New(reflect.TypeOf(list).Elem()).Interface().(LIST)
		opt := client.MatchingFields{}
		if brokerConfiguration.Labels["formance.com/stack"] != "any" {
			opt[".spec.stack"] = brokerConfiguration.Labels["formance.com/stack"]
		}

		if err := mgr.GetClient().List(ctx, list, opt); err != nil {
			return []reconcile.Request{}
		}

		return MapObjectToReconcileRequests(
			ExtractItemsFromList[LIST, OBJECT](list)...,
		)
	}
}
