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
	autoscallingv1 "k8s.io/api/autoscaling/v1"
	autoscallingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Mutator reconciles a Auth object
type LedgerMutator struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=components.formance.com,resources=ledgers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=components.formance.com,resources=ledgers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=components.formance.com,resources=ledgers/finalizers,verbs=update

func (r *LedgerMutator) Mutate(ctx context.Context, ledger *componentsv1beta2.Ledger) (*ctrl.Result, error) {

	apisv1beta2.SetProgressing(ledger)

	deployment, err := r.reconcileDeployment(ctx, ledger)
	if err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling deployment")
	}

	service, err := r.reconcileService(ctx, ledger, deployment)
	if err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling service")
	}

	if ledger.Spec.Ingress != nil {
		_, err = r.reconcileIngress(ctx, ledger, service)
		if err != nil {
			return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling service")
		}
	} else {
		err = r.Client.Delete(ctx, &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ledger.Name,
				Namespace: ledger.Namespace,
			},
		})
		if err != nil && !errors.IsNotFound(err) {
			return controllerutils.Requeue(), pkgError.Wrap(err, "Deleting ingress")
		}
		apisv1beta2.RemoveIngressCondition(ledger)
	}

	if _, err := r.reconcileHPA(ctx, ledger); err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling HPA")
	}

	apisv1beta2.SetReady(ledger)

	return nil, nil
}

func (r *LedgerMutator) reconcileDeployment(ctx context.Context, ledger *componentsv1beta2.Ledger) (*appsv1.Deployment, error) {
	matchLabels := CreateMap("app.kubernetes.io/name", "ledger")

	env := []corev1.EnvVar{
		apisv1beta2.Env("NUMARY_SERVER_HTTP_BIND_ADDRESS", "0.0.0.0:8080"),
		apisv1beta2.Env("NUMARY_STORAGE_DRIVER", "postgres"),
	}
	env = append(env, ledger.Spec.Postgres.Env("NUMARY_")...)
	env = append(env, ledger.Spec.LockingStrategy.Env("NUMARY_")...)
	env = append(env, apisv1beta2.Env("NUMARY_STORAGE_POSTGRES_CONN_STRING", "$(NUMARY_POSTGRES_URI)"))
	env = append(env, ledger.Spec.DevProperties.EnvWithPrefix("NUMARY_")...)
	if ledger.Spec.Monitoring != nil {
		env = append(env, ledger.Spec.Monitoring.Env("NUMARY_")...)
	}
	if ledger.Spec.Collector != nil {
		env = append(env, ledger.Spec.Collector.Env("NUMARY_")...)
	}

	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(ledger),
		controllerutils.WithController[*appsv1.Deployment](ledger, r.Scheme),
		func(deployment *appsv1.Deployment) error {
			deployment.Spec = appsv1.DeploymentSpec{
				Replicas: ledger.Spec.GetReplicas(),
				Selector: &metav1.LabelSelector{
					MatchLabels: matchLabels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: matchLabels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:            "ledger",
							Image:           controllerutils.GetImage("ledger", ledger.Spec.Version),
							ImagePullPolicy: controllerutils.ImagePullPolicy(ledger.Spec),
							Env:             env,
							Ports: []corev1.ContainerPort{{
								Name:          "ledger",
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
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewMilliQuantity(256, resource.DecimalSI),
								},
							},
						}},
					},
				},
			}
			if ledger.Spec.Postgres.CreateDatabase {
				deployment.Spec.Template.Spec.InitContainers = []corev1.Container{
					ledger.Spec.Postgres.CreateDatabaseInitContainer(),
				}
			}
			return nil
		})
	switch {
	case err != nil:
		apisv1beta2.SetDeploymentError(ledger, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetDeploymentReady(ledger)
	}

	selector, err := metav1.LabelSelectorAsSelector(ret.Spec.Selector)
	if err != nil {
		return nil, err
	}

	ledger.Status.Selector = selector.String()
	ledger.Status.Replicas = *ledger.Spec.GetReplicas()

	return ret, err
}

func (r *LedgerMutator) reconcileHPA(ctx context.Context, ledger *componentsv1beta2.Ledger) (*autoscallingv2.HorizontalPodAutoscaler, error) {
	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(ledger),
		controllerutils.WithController[*autoscallingv2.HorizontalPodAutoscaler](ledger, r.Scheme),
		func(hpa *autoscallingv2.HorizontalPodAutoscaler) error {
			hpa.Spec = ledger.Spec.GetHPASpec(ledger)
			return nil
		})
	switch {
	case err != nil:
		apisv1beta2.SetHPAError(ledger, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetHPAReady(ledger)
	}
	return ret, err
}

func (r *LedgerMutator) reconcileService(ctx context.Context, ledger *componentsv1beta2.Ledger, deployment *appsv1.Deployment) (*corev1.Service, error) {
	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(ledger),
		controllerutils.WithController[*corev1.Service](ledger, r.Scheme),
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
		apisv1beta2.SetServiceError(ledger, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetServiceReady(ledger)
	}
	return ret, err
}

func (r *LedgerMutator) reconcileIngress(ctx context.Context, ledger *componentsv1beta2.Ledger, service *corev1.Service) (*networkingv1.Ingress, error) {
	annotations := ledger.Spec.Ingress.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	middlewareAuth := fmt.Sprintf("%s-auth-middleware@kubernetescrd", ledger.Namespace)
	annotations["traefik.ingress.kubernetes.io/router.middlewares"] = fmt.Sprintf("%s, %s", middlewareAuth, annotations["traefik.ingress.kubernetes.io/router.middlewares"])
	ret, operationResult, err := controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(ledger),
		controllerutils.WithController[*networkingv1.Ingress](ledger, r.Scheme),
		func(ingress *networkingv1.Ingress) error {
			pathType := networkingv1.PathTypePrefix
			ingress.ObjectMeta.Annotations = annotations
			ingress.Spec = networkingv1.IngressSpec{
				TLS: ledger.Spec.Ingress.TLS.AsK8SIngressTLSSlice(),
				Rules: []networkingv1.IngressRule{
					{
						Host: ledger.Spec.Ingress.Host,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     ledger.Spec.Ingress.Path,
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
		apisv1beta2.SetIngressError(ledger, err.Error())
		return nil, err
	case operationResult == controllerutil.OperationResultNone:
	default:
		apisv1beta2.SetIngressReady(ledger)
	}
	return ret, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LedgerMutator) SetupWithBuilder(mgr ctrl.Manager, builder *ctrl.Builder) error {
	builder.
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&autoscallingv1.HorizontalPodAutoscaler{})
	return nil
}

func NewLedgerMutator(client client.Client, scheme *runtime.Scheme) controllerutils.Mutator[*componentsv1beta2.Ledger] {
	return &LedgerMutator{
		Client: client,
		Scheme: scheme,
	}
}
