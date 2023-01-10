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

// CounterpartiesMutator reconciles a Auth object
type CounterpartiesMutator struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=components.formance.com,resources=counterparties,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=components.formance.com,resources=counterparties/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=components.formance.com,resources=counterparties/finalizers,verbs=update

func (r *CounterpartiesMutator) Mutate(ctx context.Context, counterparties *componentsv1beta2.Counterparties) (*ctrl.Result, error) {

	apisv1beta2.SetProgressing(counterparties)

	if counterparties.Spec.Enabled {
		deployment, err := r.reconcileDeployment(ctx, counterparties)
		if err != nil {
			return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling deployment")
		}

		service, err := r.reconcileService(ctx, counterparties, deployment)
		if err != nil {
			return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling service")
		}

		if counterparties.Spec.Ingress != nil {
			_, err = r.reconcileIngress(ctx, counterparties, service)
			if err != nil {
				return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling service")
			}
		} else {
			err = r.Client.Delete(ctx, &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      counterparties.Name,
					Namespace: counterparties.Namespace,
				},
			})
			if err != nil && !errors.IsNotFound(err) {
				return controllerutils.Requeue(), pkgError.Wrap(err, "Deleting ingress")
			}
			apisv1beta2.RemoveIngressCondition(counterparties)
		}
	}

	apisv1beta2.SetReady(counterparties)

	return nil, nil
}

func counterpartiesEnvVars(counterparties *componentsv1beta2.Counterparties) []corev1.EnvVar {
	env := counterparties.Spec.Postgres.Env("")

	env = append(env, counterparties.Spec.DevProperties.Env()...)
	if counterparties.Spec.Monitoring != nil {
		env = append(env, counterparties.Spec.Monitoring.Env("")...)
	}
	return env
}

func (r *CounterpartiesMutator) reconcileDeployment(ctx context.Context, counterparties *componentsv1beta2.Counterparties) (*appsv1.Deployment, error) {
	matchLabels := CreateMap("app.kubernetes.io/name", "counterparties")

	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(counterparties),
		controllerutils.WithController[*appsv1.Deployment](counterparties, r.Scheme),
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
							Name:            "counterparties",
							Image:           controllerutils.GetImage("counterparties", counterparties.Spec.Version),
							ImagePullPolicy: controllerutils.ImagePullPolicy(counterparties.Spec),
							Env:             counterpartiesEnvVars(counterparties),
							Ports: []corev1.ContainerPort{{
								Name:          "counterparties",
								ContainerPort: 8080,
							}},
							LivenessProbe: controllerutils.DefaultLiveness(),
						}},
					},
				},
			}
			if counterparties.Spec.Postgres.CreateDatabase {
				deployment.Spec.Template.Spec.InitContainers = []corev1.Container{{
					Name:            "init-create-counterparties-db",
					Image:           "postgres:13",
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command: []string{
						"sh",
						"-c",
						`psql -Atx ${POSTGRES_NO_DATABASE_URI}/postgres -c "SELECT 1 FROM pg_database WHERE datname = '${POSTGRES_DATABASE}'" | grep -q 1 && echo "Base already exists" || psql -Atx ${POSTGRES_NO_DATABASE_URI}/postgres -c "CREATE DATABASE \"${POSTGRES_DATABASE}\""`,
					},
					Env: counterparties.Spec.Postgres.Env(""),
				}}
			}
			return nil
		})
	switch {
	case err != nil:
		apisv1beta2.SetDeploymentError(counterparties, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetDeploymentReady(counterparties)
	}
	return ret, err
}

func (r *CounterpartiesMutator) reconcileService(ctx context.Context, counterparties *componentsv1beta2.Counterparties, deployment *appsv1.Deployment) (*corev1.Service, error) {
	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(counterparties),
		controllerutils.WithController[*corev1.Service](counterparties, r.Scheme),
		func(service *corev1.Service) error {
			service.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:        "counterparties",
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
		apisv1beta2.SetServiceError(counterparties, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetServiceReady(counterparties)
	}
	return ret, err
}

func (r *CounterpartiesMutator) reconcileIngress(ctx context.Context, counterparties *componentsv1beta2.Counterparties, service *corev1.Service) (*networkingv1.Ingress, error) {
	annotations := counterparties.Spec.Ingress.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	middlewareAuth := fmt.Sprintf("%s-auth-middleware@kubernetescrd", counterparties.Namespace)
	annotations["traefik.ingress.kubernetes.io/router.middlewares"] = fmt.Sprintf("%s, %s", middlewareAuth, annotations["traefik.ingress.kubernetes.io/router.middlewares"])
	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(counterparties),
		controllerutils.WithController[*networkingv1.Ingress](counterparties, r.Scheme),
		func(ingress *networkingv1.Ingress) error {
			pathType := networkingv1.PathTypePrefix
			ingress.ObjectMeta.Annotations = annotations
			ingress.Spec = networkingv1.IngressSpec{
				TLS: counterparties.Spec.Ingress.TLS.AsK8SIngressTLSSlice(),
				Rules: []networkingv1.IngressRule{
					{
						Host: counterparties.Spec.Ingress.Host,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     counterparties.Spec.Ingress.Path,
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
		apisv1beta2.SetIngressError(counterparties, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetIngressReady(counterparties)
	}
	return ret, nil
}

// SetupWithBuilder SetupWithManager sets up the controller with the Manager.
func (r *CounterpartiesMutator) SetupWithBuilder(mgr ctrl.Manager, builder *ctrl.Builder) error {
	builder.
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{})
	return nil
}

func NewCounterpartiesMutator(client client.Client, scheme *runtime.Scheme) controllerutils.Mutator[*componentsv1beta2.Counterparties] {
	return &CounterpartiesMutator{
		Client: client,
		Scheme: scheme,
	}
}
