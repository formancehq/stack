package databases

import (
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Watch(service string, target any) func(ctx core.Context, object *v1beta1.Database) []reconcile.Request {
	return func(ctx core.Context, object *v1beta1.Database) []reconcile.Request {
		if object.Spec.Service != service {
			return []reconcile.Request{}
		}

		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(target)), 0, 0).Interface()

		err := core.GetAllDependents(ctx, object.Spec.Stack, &slice)
		if err != nil {
			return []reconcile.Request{}
		}

		objects := make([]client.Object, 0)
		for i := 0; i < reflect.ValueOf(slice).Len(); i++ {
			objects = append(objects, reflect.ValueOf(slice).Index(i).Interface().(client.Object))
		}

		return core.MapObjectToReconcileRequests(objects...)
	}
}
