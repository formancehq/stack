package core

import (
	"context"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func WatchDependents(mgr Manager, t client.Object) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {

		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(t)), 0, 0).Interface()

		err := GetAllStackDependencies(
			NewContext(mgr, ctx),
			object.(v1beta1.Dependent).GetStack(), &slice)
		if err != nil {
			return nil
		}

		objects := make([]client.Object, 0)
		for i := 0; i < reflect.ValueOf(slice).Len(); i++ {
			objects = append(objects, reflect.ValueOf(slice).Index(i).Interface().(client.Object))
		}

		return MapObjectToReconcileRequests(objects...)
	}
}

func Watch(mgr Manager, t client.Object) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {

		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(t)), 0, 0).Interface()

		err := GetAllStackDependencies(
			NewContext(mgr, ctx),
			object.(*v1beta1.Stack).Name, &slice)
		if err != nil {
			return nil
		}

		objects := make([]client.Object, 0)
		for i := 0; i < reflect.ValueOf(slice).Len(); i++ {
			objects = append(objects, reflect.ValueOf(slice).Index(i).Interface().(client.Object))
		}

		return MapObjectToReconcileRequests(objects...)
	}
}
