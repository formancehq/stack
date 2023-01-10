/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package components

import (
	"context"
	"fmt"

	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
	apisv1beta2 "github.com/formancehq/operator/pkg/apis/v1beta2"
	"github.com/formancehq/operator/pkg/controllerutils"
	. "github.com/formancehq/operator/pkg/typeutils"
	pkgError "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Mutator reconciles a Auth object
type PaymentsMutator struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=components.formance.com,resources=payments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=components.formance.com,resources=payments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=components.formance.com,resources=payments/finalizers,verbs=update

func (r *PaymentsMutator) Mutate(ctx context.Context, payments *componentsv1beta2.Payments) (*ctrl.Result, error) {

	apisv1beta2.SetProgressing(payments)

	deployment, err := r.reconcileDeployment(ctx, payments)
	if err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling deployment")
	}

	service, err := r.reconcileService(ctx, payments, deployment)
	if err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling service")
	}

	if payments.Spec.Ingress != nil {
		_, err = r.reconcileIngress(ctx, payments, service)
		if err != nil {
			return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling service")
		}
	} else {
		err = r.Client.Delete(ctx, &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      payments.Name,
				Namespace: payments.Namespace,
			},
		})
		if err != nil && !errors.IsNotFound(err) {
			return controllerutils.Requeue(), pkgError.Wrap(err, "Deleting ingress")
		}
		apisv1beta2.RemoveIngressCondition(payments)
	}

	apisv1beta2.SetReady(payments)

	return nil, nil
}

func (r *PaymentsMutator) reconcileDeployment(ctx context.Context, payments *componentsv1beta2.Payments) (*appsv1.Deployment, error) {
	matchLabels := CreateMap("app.kubernetes.io/name", "payments")

	env := payments.Spec.Postgres.Env("")
	env = append(env,
		apisv1beta2.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
	)
	if payments.Spec.Debug {
		env = append(env, apisv1beta2.Env("DEBUG", "true"))
	}
	if payments.Spec.Monitoring != nil {
		env = append(env, payments.Spec.Monitoring.Env("")...)
	}
	if payments.Spec.Collector != nil {
		env = append(env, payments.Spec.Collector.Env("")...)
	}

	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(payments),
		controllerutils.WithController[*appsv1.Deployment](payments, r.Scheme),
		func(deployment *appsv1.Deployment) error {
			deployment.Spec = appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: matchLabels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: matchLabels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:            "payments",
							Image:           controllerutils.GetImage("payments", payments.Spec.Version),
							ImagePullPolicy: controllerutils.ImagePullPolicy(payments.Spec),
							Env:             env,
							Ports: []corev1.ContainerPort{{
								Name:          "payments",
								ContainerPort: 8080,
							}},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/_health",
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
										Scheme: "HTTP",
									},
								},
								InitialDelaySeconds:           1,
								TimeoutSeconds:                30,
								PeriodSeconds:                 2,
								SuccessThreshold:              1,
								FailureThreshold:              10,
								TerminationGracePeriodSeconds: pointer.Int64(10),
							},
						}},
					},
				},
			}
			if payments.Spec.Postgres.CreateDatabase {
				deployment.Spec.Template.Spec.InitContainers = []corev1.Container{
					payments.Spec.Postgres.CreateDatabaseInitContainer(),
					{
						Name:            "migrate",
						Image:           controllerutils.GetImage("payments", payments.Spec.Version),
						ImagePullPolicy: controllerutils.ImagePullPolicy(payments.Spec),
						Env:             env,
						Command:         []string{"payments", "migrate", "up"},
					}}
			}
			return nil
		})
	switch {
	case err != nil:
		apisv1beta2.SetDeploymentError(payments, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetDeploymentReady(payments)
	}
	return ret, err
}

func (r *PaymentsMutator) reconcileService(ctx context.Context, payments *componentsv1beta2.Payments, deployment *appsv1.Deployment) (*corev1.Service, error) {
	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(payments),
		controllerutils.WithController[*corev1.Service](payments, r.Scheme),
		func(service *corev1.Service) error {
			service.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:        "http",
					Port:        deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort,
					Protocol:    "TCP",
					AppProtocol: pointer.String("http"),
					TargetPort:  intstr.FromString(deployment.Spec.Template.Spec.Containers[0].Ports[0].Name),
				}},
				Selector: deployment.Spec.Template.Labels,
			}
			return nil
		})
	switch {
	case err != nil:
		apisv1beta2.SetServiceError(payments, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetServiceReady(payments)
	}
	return ret, err
}

func (r *PaymentsMutator) reconcileIngress(ctx context.Context, payments *componentsv1beta2.Payments, service *corev1.Service) (*networkingv1.Ingress, error) {
	annotations := payments.Spec.Ingress.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	middlewareAuth := fmt.Sprintf("%s-auth-middleware@kubernetescrd", payments.Namespace)
	annotations["traefik.ingress.kubernetes.io/router.middlewares"] = fmt.Sprintf("%s, %s", middlewareAuth, annotations["traefik.ingress.kubernetes.io/router.middlewares"])
	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(payments),
		controllerutils.WithController[*networkingv1.Ingress](payments, r.Scheme),
		func(ingress *networkingv1.Ingress) error {
			pathType := networkingv1.PathTypePrefix
			ingress.ObjectMeta.Annotations = annotations
			ingress.Spec = networkingv1.IngressSpec{
				TLS: payments.Spec.Ingress.TLS.AsK8SIngressTLSSlice(),
				Rules: []networkingv1.IngressRule{
					{
						Host: payments.Spec.Ingress.Host,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     payments.Spec.Ingress.Path,
										PathType: &pathType,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: service.Name,
												Port: networkingv1.ServiceBackendPort{
													Name: service.Spec.Ports[0].Name,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			return nil
		})
	switch {
	case err != nil:
		apisv1beta2.SetIngressError(payments, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetIngressReady(payments)
	}
	return ret, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PaymentsMutator) SetupWithBuilder(mgr ctrl.Manager, builder *ctrl.Builder) error {
	builder.
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{})
	return nil
}

func NewPaymentsMutator(client client.Client, scheme *runtime.Scheme) controllerutils.Mutator[*componentsv1beta2.Payments] {
	return &PaymentsMutator{
		Client: client,
		Scheme: scheme,
	}
}
