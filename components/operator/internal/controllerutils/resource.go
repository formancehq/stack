package controllerutils

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/api/equality"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ObjectFactory[T client.Object] struct {
	client    client.Client
	options   []ObjectMutator[T]
	namespace string
}

func (f *ObjectFactory[T]) CreateOrUpdate(ctx context.Context, name string, options ...ObjectMutator[T]) (T, error) {
	ret, _, err := CreateOrUpdate[T](ctx, f.client, types.NamespacedName{
		Namespace: f.namespace,
		Name:      name,
	}, append(options, f.options...)...)
	if err != nil {
		var ret T
		return ret, err
	}
	return ret, nil
}

func (f *ObjectFactory[T]) Update(ctx context.Context, t T, options ...ObjectMutator[T]) (T, error) {
	ret, _, err := Update[T](ctx, f.client, t, append(options, f.options...)...)
	if err != nil {
		var ret T
		return ret, err
	}
	return ret, nil
}

func (f *ObjectFactory[T]) Get(ctx context.Context, name string) (T, error) {
	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	if err := f.client.Get(ctx, types.NamespacedName{
		Namespace: f.namespace,
		Name:      name,
	}, t); err != nil {
		return t, err
	}
	return t, nil
}

func NewObjectFactory[T client.Object](client client.Client, ns string, options ...ObjectMutator[T]) *ObjectFactory[T] {
	return &ObjectFactory[T]{
		namespace: ns,
		client:    client,
		options:   options,
	}
}

type ObjectMutator[T any] func(t T)

func WithController[T client.Object](owner client.Object, scheme *runtime.Scheme) ObjectMutator[T] {
	return func(t T) {
		if !v1.IsControlledBy(t, owner) {
			if err := controllerutil.SetControllerReference(owner, t, scheme); err != nil {
				panic(err)
			}
		}
	}
}

func WithAnnotations[T client.Object](newAnnotations map[string]string) ObjectMutator[T] {
	return func(t T) {
		annotations := t.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}
		for k, v := range newAnnotations {
			annotations[k] = v
		}
		t.SetAnnotations(annotations)
	}
}

func CreateOrUpdate[T client.Object](ctx context.Context, client client.Client,
	key types.NamespacedName, mutators ...ObjectMutator[T]) (T, controllerutil.OperationResult, error) {

	var ret T
	ret = reflect.New(reflect.TypeOf(ret).Elem()).Interface().(T)
	ret.SetNamespace(key.Namespace)
	ret.SetName(key.Name)
	operationResult, err := controllerutil.CreateOrUpdate(ctx, client, ret, func() error {
		labels := ret.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}
		labels["stack"] = "true"
		ret.SetLabels(labels)
		for _, mutate := range mutators {
			mutate(ret)
		}
		return nil
	})
	return ret, operationResult, err
}

func Update[T client.Object](ctx context.Context, client client.Client,
	t T, mutators ...ObjectMutator[T]) (T, controllerutil.OperationResult, error) {

	var ret T
	existing := t.DeepCopyObject().(T)
	for _, mutate := range mutators {
		mutate(existing)
	}

	if equality.Semantic.DeepEqual(existing, t) {
		return t, controllerutil.OperationResultNone, nil
	}

	if err := client.Update(ctx, existing); err != nil {
		return ret, controllerutil.OperationResultNone, err
	}
	return ret, controllerutil.OperationResultUpdated, nil
}
