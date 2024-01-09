package stacks

import (
	"context"
	"github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func GetStack(ctx core.Context, spec Dependent) (*v1beta1.Stack, error) {
	stack := &v1beta1.Stack{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: spec.GetStack(),
	}, stack); err != nil {
		return nil, err
	}

	return stack, nil
}

type Dependent interface {
	GetStack() string
}

var (
	ErrNotFound               = errors.New("no configuration found")
	ErrMultipleInstancesFound = errors.New("multiple resources found")
)

func GetAllDependents[T client.Object](ctx core.Context, stackName string) ([]T, error) {
	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)

	kinds, _, err := ctx.GetScheme().ObjectKinds(t)
	if err != nil {
		return nil, err
	}

	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(kinds[0])

	err = ctx.GetClient().List(ctx, list, client.MatchingFields{
		"stack": stackName,
	})
	if err != nil {
		return nil, err
	}

	return collectionutils.Map(list.Items, func(from unstructured.Unstructured) T {
		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(from.Object, t); err != nil {
			panic(err)
		}
		return t
	}), nil
}

func GetSingleDependency[T client.Object](ctx core.Context, stackName string) (T, error) {
	var zeroValue T

	items, err := GetAllDependents[T](ctx, stackName)
	if err != nil {
		return zeroValue, err
	}

	switch len(items) {
	case 0:
		return zeroValue, nil
	case 1:
		return items[0], nil
	default:
		return zeroValue, ErrMultipleInstancesFound
	}
}

func HasDependency[T client.Object](ctx core.Context, stackName string) (bool, error) {
	ret, err := GetSingleDependency[T](ctx, stackName)
	if err != nil && !errors.Is(err, ErrMultipleInstancesFound) {
		return false, err
	}
	if reflect.ValueOf(ret).IsZero() {
		return false, nil
	}
	return true, nil
}

func GetIfEnabled[T client.Object](ctx core.Context, stackName string) (T, error) {
	var t T
	ret, err := GetSingleDependency[T](ctx, stackName)
	if err != nil {
		return t, err
	}
	if reflect.ValueOf(ret).IsZero() {
		return t, nil
	}

	return ret, nil
}

func IsEnabledByLabel[T client.Object](ctx core.Context, stackName string) (bool, error) {
	ret, err := GetByLabel[T](ctx, stackName)
	if err != nil {
		return false, err
	}
	if reflect.ValueOf(ret).IsZero() {
		return false, nil
	}

	return true, nil
}

func WatchUsingLabels[T client.Object](mgr core.Manager) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		options := make([]client.ListOption, 0)
		if object.GetLabels()[core.StackLabel] != "any" {
			options = append(options, client.MatchingFields{
				"stack": object.GetLabels()[core.StackLabel],
			})
		}

		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		kinds, _, err := mgr.GetScheme().ObjectKinds(t)
		if err != nil {
			return []reconcile.Request{}
		}

		us := &unstructured.UnstructuredList{}
		us.SetGroupVersionKind(kinds[0])
		if err := mgr.GetClient().List(ctx, us, options...); err != nil {
			return []reconcile.Request{}
		}

		return core.MapObjectToReconcileRequests(
			collectionutils.Map(us.Items, collectionutils.ToPointer[unstructured.Unstructured])...,
		)
	}
}

func GetByLabel[T client.Object](ctx core.Context, stackName string) (T, error) {

	var (
		zeroValue T
	)
	stackSelectorRequirement, err := labels.NewRequirement(core.StackLabel, selection.In, []string{"any", stackName})
	if err != nil {
		return zeroValue, err
	}

	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	kinds, _, err := ctx.GetScheme().ObjectKinds(t)
	if err != nil {
		return zeroValue, err
	}

	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(kinds[0])

	if err := ctx.GetClient().List(ctx, list, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
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
		return zeroValue, errors.New("found multiple broker config")
	}
}

func Require[T client.Object](ctx core.Context, stackName string) (T, error) {
	t, err := GetByLabel[T](ctx, stackName)
	if err != nil {
		return t, err
	}
	if reflect.ValueOf(t).IsZero() {
		return t, errors.New("not found")
	}
	return t, nil
}
