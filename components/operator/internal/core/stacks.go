package core

import (
	"context"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func GetStack(ctx Context, spec Dependent) (*v1beta1.Stack, error) {
	stack := &v1beta1.Stack{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: spec.GetStack(),
	}, stack); err != nil {
		return nil, err
	}

	return stack, nil
}

type Object interface {
	client.Object
	SetReady(bool)
	IsReady() bool
	SetError(string)
}

type Dependent interface {
	Object
	GetStack() string
}

var (
	ErrNotFound               = errors.New("no configuration found")
	ErrMultipleInstancesFound = errors.New("multiple resources found")
)

func GetAllDependents(ctx Context, stackName string, to any) error {
	slice := reflect.Indirect(reflect.ValueOf(to)).Interface()
	objectType := reflect.TypeOf(slice).Elem()

	kinds, _, err := ctx.GetScheme().ObjectKinds(reflect.New(objectType.Elem()).Interface().(client.Object))
	if err != nil {
		return err
	}

	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(kinds[0])

	err = ctx.GetClient().List(ctx, list, client.MatchingFields{
		"stack": stackName,
	})
	if err != nil {
		return err
	}

	ret := reflect.ValueOf(slice)
	for _, item := range list.Items {
		t := reflect.New(objectType.Elem()).Interface().(client.Object)
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, t); err != nil {
			panic(err)
		}
		ret = reflect.Append(ret, reflect.ValueOf(t))
	}

	reflect.ValueOf(to).Elem().Set(ret)

	return nil
}

func GetSingleDependency(ctx Context, stackName string, to client.Object) error {

	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(to)), 0, 0).Interface()
	err := GetAllDependents(ctx, stackName, &slice)
	if err != nil {
		return err
	}

	switch reflect.ValueOf(slice).Len() {
	case 0:
		return ErrNotFound
	case 1:
		reflect.Indirect(reflect.ValueOf(to)).Set(reflect.ValueOf(slice).Index(0).Elem())
		return nil
	default:
		return ErrMultipleInstancesFound
	}
}

func HasDependency(ctx Context, stackName string, to client.Object) (bool, error) {
	err := GetSingleDependency(ctx, stackName, to)
	if err != nil && !errors.Is(err, ErrMultipleInstancesFound) {
		switch {
		default:
			return false, err
		case errors.Is(err, ErrNotFound):
			return false, nil
		}
	}
	return true, nil
}

func GetIfEnabled(ctx Context, stackName string, to client.Object) (bool, error) {
	err := GetSingleDependency(ctx, stackName, to)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return false, err
	}
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	return true, nil
}

func IsEnabledByLabel[T client.Object](ctx Context, stackName string) (bool, error) {
	ret, err := GetByLabel[T](ctx, stackName)
	if err != nil {
		return false, err
	}
	if reflect.ValueOf(ret).IsZero() {
		return false, nil
	}

	return true, nil
}

func WatchUsingLabels(mgr Manager, t client.Object) func(ctx context.Context, object client.Object) []reconcile.Request {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		options := make([]client.ListOption, 0)
		if object.GetLabels()[StackLabel] != "any" {
			options = append(options, client.MatchingFields{
				"stack": object.GetLabels()[StackLabel],
			})
		}

		kinds, _, err := mgr.GetScheme().ObjectKinds(t)
		if err != nil {
			return []reconcile.Request{}
		}

		us := &unstructured.UnstructuredList{}
		us.SetGroupVersionKind(kinds[0])
		if err := mgr.GetClient().List(ctx, us, options...); err != nil {
			return []reconcile.Request{}
		}

		return MapObjectToReconcileRequests(
			collectionutils.Map(us.Items, collectionutils.ToPointer[unstructured.Unstructured])...,
		)
	}
}

func GetByLabel[T client.Object](ctx Context, stackName string) (T, error) {

	var (
		zeroValue T
	)
	stackSelectorRequirement, err := labels.NewRequirement(StackLabel, selection.In, []string{"any", stackName})
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
		return zeroValue, errors.New("found multiple configuration")
	}
}

func RequireLabelledConfig[T client.Object](ctx Context, stackName string) (T, error) {
	t, err := GetByLabel[T](ctx, stackName)
	if err != nil {
		return t, err
	}
	if reflect.ValueOf(t).IsZero() {
		return t, errors.New("not found")
	}
	return t, nil
}
