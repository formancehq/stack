package databases

import (
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch[T client.Object]() core.ReconcilerOption[T] {
	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	return core.WithWatch[T, *v1beta1.Database](func(ctx core.Context, database *v1beta1.Database) []reconcile.Request {
		if database.Spec.Service != core.LowerCamelCaseKind(ctx, t) {
			return []reconcile.Request{}
		}

		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(t)), 0, 0).Interface()

		err := core.GetAllStackDependencies(ctx, database.Spec.Stack, &slice)
		if err != nil {
			return []reconcile.Request{}
		}

		objects := make([]client.Object, 0)
		for i := 0; i < reflect.ValueOf(slice).Len(); i++ {
			objects = append(objects, reflect.ValueOf(slice).Index(i).Interface().(client.Object))
		}

		return core.MapObjectToReconcileRequests(objects...)
	})
}
