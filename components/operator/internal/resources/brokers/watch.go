package brokers

import (
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch[T client.Object]() core.ReconcilerOption[T] {
	return core.WithWatch[T, *v1beta1.Broker](func(ctx core.Context, topic *v1beta1.Broker) []reconcile.Request {
		var t T
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(t)), 0, 0).Interface()

		err := core.GetAllStackDependencies(ctx, topic.Spec.Stack, &slice)
		if err != nil {

			return nil
		}

		objects := make([]client.Object, 0)
		for i := 0; i < reflect.ValueOf(slice).Len(); i++ {
			objects = append(objects, reflect.ValueOf(slice).Index(i).Interface().(client.Object))
		}

		return core.MapObjectToReconcileRequests(objects...)
	})
}
