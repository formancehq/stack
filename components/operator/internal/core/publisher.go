package core

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ForEachEventPublisher(ctx Context, stackName string, fn func(object client.Object) error) error {
	ret, err := ListEventPublishers(ctx, stackName)
	if err != nil {
		return err
	}

	for _, object := range ret {
		if err := fn(&object); err != nil {
			return err
		}
	}

	return nil
}

func ListEventPublishers(ctx Context, stackName string) ([]unstructured.Unstructured, error) {
	ret := make([]unstructured.Unstructured, 0)
	for gvk, rtype := range ctx.GetScheme().AllKnownTypes() {
		object, ok := reflect.New(rtype).Interface().(client.Object)
		if !ok {
			continue
		}

		if _, ok := object.(v1beta1.EventPublisher); ok {
			us := &unstructured.UnstructuredList{}
			us.SetGroupVersionKind(gvk)

			if err := ctx.GetClient().List(ctx, us, client.MatchingFields{
				"stack": stackName,
			}); err != nil {
				return nil, err
			}

			for _, item := range us.Items {
				item.SetGroupVersionKind(gvk)
				ret = append(ret, item)
			}
		}
	}

	return ret, nil
}
