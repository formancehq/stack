package stacks

import (
	errors2 "errors"
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

func GetDependentObjects[LIST client.ObjectList, OBJECT client.Object](ctx core.Context, stackName string) ([]OBJECT, error) {
	var list LIST
	list = reflect.New(reflect.TypeOf(list).Elem()).Interface().(LIST)

	err := ctx.GetClient().List(ctx, list, client.MatchingFields{
		".spec.stack": stackName,
	})
	if err != nil {
		return nil, err
	}

	items := reflect.ValueOf(list).
		Elem().
		FieldByName("Items")

	ret := make([]OBJECT, 0)
	for i := 0; i < items.Len(); i++ {
		ret = append(ret, items.Index(i).Addr().Interface().(OBJECT))
	}

	return ret, nil
}

func GetSingleStackDependencyObject[LIST client.ObjectList, OBJECT client.Object](ctx core.Context, stackName string) (OBJECT, error) {

	var t OBJECT

	items, err := GetDependentObjects[LIST, OBJECT](ctx, stackName)
	if err != nil {
		return t, err
	}

	switch len(items) {
	case 0:
		return t, nil
	case 1:
		return items[0], nil
	default:
		return t, ErrMultipleInstancesFound
	}
}

func HasSingleStackDependencyObject[LIST client.ObjectList, OBJECT client.Object](ctx core.Context, stackName string) (bool, error) {
	ret, err := GetSingleStackDependencyObject[LIST, OBJECT](ctx, stackName)
	if err != nil && !errors2.Is(err, ErrMultipleInstancesFound) {
		return false, err
	}
	if reflect.ValueOf(ret).IsZero() {
		return false, nil
	}
	return true, nil
}

func RequireSingleStackDependencyObject[LIST client.ObjectList, OBJECT client.Object](ctx core.Context, stackName string) (OBJECT, error) {
	var ret OBJECT
	ret, err := GetSingleStackDependencyObject[LIST, OBJECT](ctx, stackName)
	if err != nil && !errors2.Is(err, ErrMultipleInstancesFound) {
		return ret, err
	}
	if reflect.ValueOf(ret).Elem().IsZero() {
		return ret, ErrNotFound
	}
	return ret, nil
}
