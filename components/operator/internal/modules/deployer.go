package modules

import (
	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/common"
	"github.com/formancehq/operator/internal/controllerutils"
	traefik "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
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
	Middlewares(...controllerutils.ObjectMutator[*traefik.Middleware]) *controllerutils.ObjectFactory[*traefik.Middleware]
	Services(...controllerutils.ObjectMutator[*corev1.Service]) *controllerutils.ObjectFactory[*corev1.Service]
	Ingresses(...controllerutils.ObjectMutator[*networkingv1.Ingress]) *controllerutils.ObjectFactory[*networkingv1.Ingress]
}

type StackWideDeployer struct {
	client        client.Client
	scheme        *runtime.Scheme
	stack         *v1beta3.Stack
	configuration *v1beta3.Configuration
}

func (d *StackWideDeployer) Ingresses(options ...controllerutils.ObjectMutator[*networkingv1.Ingress]) *controllerutils.ObjectFactory[*networkingv1.Ingress] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*networkingv1.Ingress](d.stack, d.scheme)...,
	)...)
}

func (d *StackWideDeployer) Services(options ...controllerutils.ObjectMutator[*corev1.Service]) *controllerutils.ObjectFactory[*corev1.Service] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.Service](d.stack, d.scheme)...,
	)...)
}

func (d *StackWideDeployer) Deployments(options ...controllerutils.ObjectMutator[*appsv1.Deployment]) *controllerutils.ObjectFactory[*appsv1.Deployment] {
	options = append(options,
		CommonOptions[*appsv1.Deployment](d.stack, d.scheme)...,
	)
	options = append(options, common.WithReloaderAnnotations[*appsv1.Deployment]())
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, options...)
}

func (d *StackWideDeployer) ConfigMaps(options ...controllerutils.ObjectMutator[*corev1.ConfigMap]) *controllerutils.ObjectFactory[*corev1.ConfigMap] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.ConfigMap](d.stack, d.scheme)...,
	)...)
}

func (d *StackWideDeployer) Secrets(options ...controllerutils.ObjectMutator[*corev1.Secret]) *controllerutils.ObjectFactory[*corev1.Secret] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.Secret](d.stack, d.scheme)...,
	)...)
}

func (d *StackWideDeployer) Middlewares(options ...controllerutils.ObjectMutator[*traefik.Middleware]) *controllerutils.ObjectFactory[*traefik.Middleware] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*traefik.Middleware](d.stack, d.scheme)...,
	)...)
}

func (d *StackWideDeployer) ForService(service string) *ComponentDeployer {
	return NewComponentDeployer(d.client, d.scheme, d.stack, d.configuration, service)
}

var _ Deployer = &StackWideDeployer{}

func NewDeployer(client client.Client, scheme *runtime.Scheme, stack *v1beta3.Stack,
	configuration *v1beta3.Configuration) *StackWideDeployer {
	return &StackWideDeployer{
		client:        client,
		scheme:        scheme,
		stack:         stack,
		configuration: configuration,
	}
}

type ComponentDeployer struct {
	StackWideDeployer
	serviceName string
}

var _ Deployer = &ComponentDeployer{}

func (d *ComponentDeployer) Deployments(options ...controllerutils.ObjectMutator[*appsv1.Deployment]) *controllerutils.ObjectFactory[*appsv1.Deployment] {
	return d.StackWideDeployer.Deployments(options...)
}

func NewComponentDeployer(client client.Client, scheme *runtime.Scheme, stack *v1beta3.Stack,
	configuration *v1beta3.Configuration, serviceName string) *ComponentDeployer {
	return &ComponentDeployer{
		StackWideDeployer: *NewDeployer(client, scheme, stack, configuration),
		serviceName:       serviceName,
	}
}

func CommonOptions[T client.Object](stack *v1beta3.Stack, scheme *runtime.Scheme) []controllerutils.ObjectMutator[T] {
	return []controllerutils.ObjectMutator[T]{
		controllerutils.WithController[T](stack, scheme),
	}
}
