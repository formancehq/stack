package modules

import (
	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/common"
	"github.com/formancehq/operator/internal/controllerutils"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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
	Jobs(...controllerutils.ObjectMutator[*batchv1.Job]) *controllerutils.ObjectFactory[*batchv1.Job]
}

type scopedResourceDeployer struct {
	client client.Client
	scheme *runtime.Scheme
	stack  *v1beta3.Stack
	owner  client.Object
}

func (d *scopedResourceDeployer) Ingresses(options ...controllerutils.ObjectMutator[*networkingv1.Ingress]) *controllerutils.ObjectFactory[*networkingv1.Ingress] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*networkingv1.Ingress](d.owner, d.scheme)...,
	)...)
}

func (d *scopedResourceDeployer) Services(options ...controllerutils.ObjectMutator[*corev1.Service]) *controllerutils.ObjectFactory[*corev1.Service] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.Service](d.owner, d.scheme)...,
	)...)
}

func (d *scopedResourceDeployer) Deployments(options ...controllerutils.ObjectMutator[*appsv1.Deployment]) *controllerutils.ObjectFactory[*appsv1.Deployment] {
	options = append(options,
		CommonOptions[*appsv1.Deployment](d.owner, d.scheme)...,
	)
	options = append(options, common.WithReloaderAnnotations[*appsv1.Deployment]())
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, options...)
}

func (d *scopedResourceDeployer) Migrations(options ...controllerutils.ObjectMutator[*v1beta3.Migration]) *controllerutils.ObjectFactory[*v1beta3.Migration] {
	options = append(options,
		CommonOptions[*v1beta3.Migration](d.owner, d.scheme)...,
	)
	options = append(options, common.WithReloaderAnnotations[*v1beta3.Migration]())
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, options...)
}

func (d *scopedResourceDeployer) ConfigMaps(options ...controllerutils.ObjectMutator[*corev1.ConfigMap]) *controllerutils.ObjectFactory[*corev1.ConfigMap] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.ConfigMap](d.owner, d.scheme)...,
	)...)
}

func (d *scopedResourceDeployer) Jobs(options ...controllerutils.ObjectMutator[*batchv1.Job]) *controllerutils.ObjectFactory[*batchv1.Job] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*batchv1.Job](d.owner, d.scheme)...,
	)...)
}

func (d *scopedResourceDeployer) CronJobs(options ...controllerutils.ObjectMutator[*batchv1.CronJob]) *controllerutils.ObjectFactory[*batchv1.CronJob] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*batchv1.CronJob](d.owner, d.scheme)...,
	)...)
}

func (d *scopedResourceDeployer) Secrets(options ...controllerutils.ObjectMutator[*corev1.Secret]) *controllerutils.ObjectFactory[*corev1.Secret] {
	return controllerutils.NewObjectFactory(d.client, d.stack.Name, append(options,
		CommonOptions[*corev1.Secret](d.owner, d.scheme)...,
	)...)
}

var _ Deployer = &scopedResourceDeployer{}

func NewScopedDeployer(client client.Client, scheme *runtime.Scheme, stack *v1beta3.Stack, owner client.Object) *scopedResourceDeployer {
	return &scopedResourceDeployer{
		client: client,
		scheme: scheme,
		stack:  stack,
		owner:  owner,
	}
}

func CommonOptions[T client.Object](owner client.Object, scheme *runtime.Scheme) []controllerutils.ObjectMutator[T] {
	return []controllerutils.ObjectMutator[T]{
		controllerutils.WithController[T](owner, scheme),
	}
}
