package internal

import (
	"bytes"
	"context"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"text/template"

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

func WithInitContainers(containers ...corev1.Container) func(r *appsv1.Deployment) {
	return func(r *appsv1.Deployment) {
		r.Spec.Template.Spec.InitContainers = containers
	}
}

func WithVolumes(volumes ...corev1.Volume) func(t *appsv1.Deployment) {
	return func(t *appsv1.Deployment) {
		t.Spec.Template.Spec.Volumes = volumes
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

func CreateCaddyfileConfigMap(ctx context.Context, _client client.Client, stack *v1beta1.Stack,
	name, _tpl string, additionalData map[string]any, options ...ObjectMutator[*corev1.ConfigMap]) (*corev1.ConfigMap, error) {
	caddyfile, err := ComputeCaddyfile(ctx, _client, stack, _tpl, additionalData)
	if err != nil {
		return nil, err
	}

	options = append([]ObjectMutator[*corev1.ConfigMap]{
		func(t *corev1.ConfigMap) {
			t.Data = map[string]string{
				"Caddyfile": caddyfile,
			}
		},
	}, options...)

	configMap, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, _client, types.NamespacedName{
		Namespace: stack.Name,
		Name:      name,
	},
		options...,
	)
	return configMap, err
}

func ComputeCaddyfile(ctx context.Context, _client client.Client, stack *v1beta1.Stack, _tpl string, additionalData map[string]any) (string, error) {
	tpl := template.Must(template.New("main").Parse(_tpl))
	buf := bytes.NewBufferString("")

	openTelemetryEnabled, err := IsOpenTelemetryEnabled(ctx, _client, stack.Name)
	if err != nil {
		return "", err
	}

	data := map[string]any{
		"EnableOpenTelemetry": openTelemetryEnabled,
	}
	data = collectionutils.MergeMaps(data, additionalData)

	if err := tpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
