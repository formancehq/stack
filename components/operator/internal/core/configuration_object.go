package core

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
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

func GetConfigurationObject[T v1beta1.ConfigurationObject](ctx Context, stackName string) (T, error) {

	var zeroValue T
	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	kinds, _, err := ctx.GetScheme().ObjectKinds(t)
	if err != nil {
		return zeroValue, err
	}

	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(kinds[0])

	if err := ctx.GetClient().List(ctx, list, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("stack", stackName),
	}); err != nil {
		return zeroValue, err
	}

	switch len(list.Items) {
	case 0:
		if err := ctx.GetClient().List(ctx, list, &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector("stack", "any"),
		}); err != nil {
			return zeroValue, err
		}
		switch len(list.Items) {
		case 0:
			return zeroValue, nil
		case 1:
			if err := runtime.DefaultUnstructuredConverter.
				FromUnstructured(list.Items[0].UnstructuredContent(), t); err != nil {
				return zeroValue, err
			}
			return t, nil
		default:
			return zeroValue, errors.New("found multiple configuration")
		}
	case 1:
		if err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(list.Items[0].UnstructuredContent(), t); err != nil {
			return zeroValue, err
		}
		return t, nil
	default:
		return zeroValue, errors.New("found multiple configuration")
	}
}

func RequireConfigurationObject[T v1beta1.ConfigurationObject](ctx Context, stackName string) (T, error) {
	t, err := GetConfigurationObject[T](ctx, stackName)
	if err != nil {
		return t, err
	}
	if reflect.ValueOf(t).IsZero() {
		return t, errors.New("not found")
	}
	return t, nil
}
