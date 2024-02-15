package core

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrNotFound               = errors.New("no configuration found")
	ErrMultipleInstancesFound = errors.New("multiple resources found")
)

func GetAllStackDependencies(ctx Context, stackName string, to any) error {
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
	err := GetAllStackDependencies(ctx, stackName, &slice)
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

func GetIfExists(ctx Context, stackName string, to client.Object) (bool, error) {
	err := GetSingleDependency(ctx, stackName, to)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return false, err
	}
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	return true, nil
}

func GetImageVersionForStack(ctx Context, stack *v1beta1.Stack, name string) (string, error) {
	if stack.Spec.Version != "" {
		return stack.Spec.Version, nil
	}
	if stack.Spec.VersionsFromFile == "" {
		return "latest", nil
	}
	versions := &v1beta1.Versions{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: stack.Spec.VersionsFromFile,
	}, versions); err != nil {
		return "", nil
	}

	version := versions.Spec[name]
	if version == "" {
		return "latest", nil
	}
	return version, nil
}
