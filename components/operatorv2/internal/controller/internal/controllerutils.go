package internal

import (
	"context"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	corev1 "k8s.io/api/core/v1"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ObjectMutator[T any] func(t T)

func WithController[T client.Object](scheme *runtime.Scheme, owner client.Object) ObjectMutator[T] {
	return func(t T) {
		if !metav1.IsControlledBy(t, owner) {
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

func WithAnnotation[T client.Object](key, value string) ObjectMutator[T] {
	return WithAnnotations[T](map[string]string{
		key: value,
	})
}

func WithMatchingLabels(name string) func(deployment *appsv1.Deployment) {
	return func(deployment *appsv1.Deployment) {
		matchLabels := map[string]string{
			"app.kubernetes.io/name": name,
		}
		if deployment.Spec.Selector == nil {
			deployment.Spec.Selector = &metav1.LabelSelector{}
		}
		deployment.Spec.Selector.MatchLabels = matchLabels
		deployment.Spec.Template.Labels = matchLabels
	}
}

func WithReplicas(replicas int32) func(t *appsv1.Deployment) {
	return func(t *appsv1.Deployment) {
		t.Spec.Replicas = pointer.For(replicas)
	}
}

func WithContainers(containers ...corev1.Container) func(r *appsv1.Deployment) {
	return func(r *appsv1.Deployment) {
		r.Spec.Template.Spec.Containers = containers
	}
}

func CreateOrUpdate[T client.Object](ctx context.Context, client client.Client,
	key types.NamespacedName, mutators ...ObjectMutator[T]) (T, controllerutil.OperationResult, error) {

	var ret T
	ret = reflect.New(reflect.TypeOf(ret).Elem()).Interface().(T)
	ret.SetNamespace(key.Namespace)
	ret.SetName(key.Name)
	operationResult, err := controllerutil.CreateOrUpdate(ctx, client, ret, func() error {
		for _, mutate := range mutators {
			mutate(ret)
		}
		return nil
	})
	return ret, operationResult, err
}
