package stacks

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetStack(ctx core.Context, spec interface {
	GetStack() string
}) (*v1beta1.Stack, error) {
	stack := &v1beta1.Stack{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: spec.GetStack(),
	}, stack); err != nil {
		return nil, err
	}

	return stack, nil
}

type stackDependent interface {
	GetStack() string
}

func IsStackDependent(object client.Object) (bool, error) {
	specField, ok := reflect.TypeOf(object).Elem().FieldByName("Spec")
	if !ok {
		return false, nil
	}

	if specField.Type.AssignableTo(reflect.TypeOf((*stackDependent)(nil)).Elem()) {
		return true, nil
	}

	return false, nil
}

func StackForDependent(object client.Object) (string, error) {
	specField, ok := reflect.TypeOf(object).Elem().FieldByName("Spec")
	if !ok {
		return "", errors.New("field Spec not found")
	}

	if specField.Type.AssignableTo(reflect.TypeOf((*stackDependent)(nil)).Elem()) {
		return reflect.ValueOf(object).
			Elem().
			FieldByName("Spec").
			Interface().(stackDependent).
			GetStack(), nil
	}

	return "", errors.New("not a stack dependent object")
}

var (
	ErrNotFound               = errors.New("no configuration found")
	ErrMultipleInstancesFound = errors.New("multiple resources found")
)

func GetDependentObjects(ctx core.Context, stackName string, list client.ObjectList) ([]client.Object, error) {
	err := ctx.GetClient().List(ctx, list, client.MatchingFields{
		".spec.stack": stackName,
	})
	if err != nil {
		return nil, err
	}

	return core.ExtractItemsFromList(list), nil
}

func GetSingleStackDependencyObject(ctx core.Context, stackName string, list client.ObjectList) (client.Object, error) {

	items, err := GetDependentObjects(ctx, stackName, list)
	if err != nil {
		return nil, err
	}

	switch len(items) {
	case 0:
		return nil, nil
	case 1:
		return items[0], nil
	default:
		return nil, ErrMultipleInstancesFound
	}
}

func HasSingleStackDependencyObject(ctx core.Context, stackName string, list client.ObjectList) (bool, error) {
	ret, err := GetSingleStackDependencyObject(ctx, stackName, list)
	if err != nil && !errors.Is(err, ErrMultipleInstancesFound) {
		return false, err
	}
	if ret == nil {
		return false, nil
	}
	return true, nil
}

func RequireSingleStackDependencyObject(ctx core.Context, stackName string, list client.ObjectList) (client.Object, error) {
	ret, err := GetSingleStackDependencyObject(ctx, stackName, list)
	if err != nil && !errors.Is(err, ErrMultipleInstancesFound) {
		return ret, err
	}
	if ret == nil {
		return ret, ErrNotFound
	}
	return ret, nil
}
