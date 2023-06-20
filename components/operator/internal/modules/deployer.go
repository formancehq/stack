package modules

import (
	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/common"
	"github.com/formancehq/operator/internal/controllerutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Deployer interface {
	Deployments(...controllerutils.ObjectMutator[*appsv1.Deployment]) *controllerutils.ObjectFactory[*appsv1.Deployment]
	ConfigMaps(...controllerutils.ObjectMutator[*corev1.ConfigMap]) *controllerutils.ObjectFactory[*corev1.ConfigMap]
	Secrets(...controllerutils.ObjectMutator[*corev1.Secret]) *controllerutils.ObjectFactory[*corev1.Secret]
	Services(...controllerutils.ObjectMutator[*corev1.Service]) *controllerutils.ObjectFactory[*corev1.Service]
	Ingresses(...controllerutils.ObjectMutator[*networkingv1.Ingress]) *controllerutils.ObjectFactory[*networkingv1.Ingress]
}

type ResourceDeployer struct {
	client        client.Client
	scheme        *runtime.Scheme
	stack         *v1beta3.Stack
	configuration *v1beta3.Configuration
}

func (d *ResourceDeployer) Ingresses(options ...controllerutils.ObjectMutator[*networkingv1.Ingress]) *controllerutils.ObjectFactory[*networkingv1.Ingress] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*networkingv1.Ingress](d.stack, d.scheme)...,
	)...)
}

func (d *ResourceDeployer) Services(options ...controllerutils.ObjectMutator[*corev1.Service]) *controllerutils.ObjectFactory[*corev1.Service] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.Service](d.stack, d.scheme)...,
	)...)
}

func (d *ResourceDeployer) Deployments(options ...controllerutils.ObjectMutator[*appsv1.Deployment]) *controllerutils.ObjectFactory[*appsv1.Deployment] {
	options = append(options,
		CommonOptions[*appsv1.Deployment](d.stack, d.scheme)...,
	)
	options = append(options, common.WithReloaderAnnotations[*appsv1.Deployment]())
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, options...)
}

func (d *ResourceDeployer) Migrations(options ...controllerutils.ObjectMutator[*v1beta3.Migration]) *controllerutils.ObjectFactory[*v1beta3.Migration] {
	options = append(options,
		CommonOptions[*v1beta3.Migration](d.stack, d.scheme)...,
	)
	options = append(options, common.WithReloaderAnnotations[*v1beta3.Migration]())
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, options...)
}

func (d *ResourceDeployer) ConfigMaps(options ...controllerutils.ObjectMutator[*corev1.ConfigMap]) *controllerutils.ObjectFactory[*corev1.ConfigMap] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.ConfigMap](d.stack, d.scheme)...,
	)...)
}

func (d *ResourceDeployer) Secrets(options ...controllerutils.ObjectMutator[*corev1.Secret]) *controllerutils.ObjectFactory[*corev1.Secret] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.Secret](d.stack, d.scheme)...,
	)...)
}

var _ Deployer = &ResourceDeployer{}

func NewDeployer(client client.Client, scheme *runtime.Scheme, stack *v1beta3.Stack,
	configuration *v1beta3.Configuration) *ResourceDeployer {
	return &ResourceDeployer{
		client:        client,
		scheme:        scheme,
		stack:         stack,
		configuration: configuration,
	}
}

func CommonOptions[T client.Object](stack *v1beta3.Stack, scheme *runtime.Scheme) []controllerutils.ObjectMutator[T] {
	return []controllerutils.ObjectMutator[T]{
		controllerutils.WithController[T](stack, scheme),
	}
}
