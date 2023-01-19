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
	"strings"

	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
	apisv1beta2 "github.com/formancehq/operator/pkg/apis/v1beta2"
	"github.com/formancehq/operator/pkg/controllerutils"
	. "github.com/formancehq/operator/pkg/typeutils"
	pkgError "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// OrchestrationMutator reconciles a Auth object
type OrchestrationMutator struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=components.formance.com,resources=orchestrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=components.formance.com,resources=orchestrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=components.formance.com,resources=orchestrations/finalizers,verbs=update

func (r *OrchestrationMutator) Mutate(ctx context.Context, orchestration *componentsv1beta2.Orchestration) (*ctrl.Result, error) {

	apisv1beta2.SetProgressing(orchestration)

	deployment, _, err := r.reconcileMainDeployment(ctx, orchestration)
	if err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling deployment")
	}

	_, _, err = r.reconcileWorkerDeployment(ctx, orchestration)
	if err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling worker deployment")
	}

	service, _, err := r.reconcileService(ctx, orchestration, deployment)
	if err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling service")
	}

	_, _, err = r.reconcileIngress(ctx, orchestration, service)
	if err != nil {
		return controllerutils.Requeue(), pkgError.Wrap(err, "Reconciling service")
	}

	apisv1beta2.SetReady(orchestration)

	return nil, nil
}

func orchestrationEnvVars(orchestration *componentsv1beta2.Orchestration) []corev1.EnvVar {
	ledgerName := strings.Replace(orchestration.GetName(), "-next", "-ledger", -1)
	env := make([]corev1.EnvVar, 0)
	env = append(env, orchestration.Spec.Postgres.Env("")...)
	env = append(env,
		apisv1beta2.Env("POSTGRES_DSN", "$(POSTGRES_URI)"),
		apisv1beta2.Env("LEDGER_URI", fmt.Sprintf("http://%s:8080", ledgerName)),
		apisv1beta2.Env("STACK_CLIENT_ID", orchestration.Spec.Auth.ClientID),
		apisv1beta2.Env("STACK_CLIENT_SECRET", orchestration.Spec.Auth.ClientSecret),
		apisv1beta2.Env("STACK_URL", orchestration.Spec.StackURL),
	)
	env = append(env, orchestration.Spec.Temporal.Env()...)

	env = append(env, orchestration.Spec.DevProperties.Env()...)
	if orchestration.Spec.Monitoring != nil {
		env = append(env, orchestration.Spec.Monitoring.Env("")...)
	}
	return env
}

func (r *OrchestrationMutator) reconcileMainDeployment(ctx context.Context, orchestration *componentsv1beta2.Orchestration) (*appsv1.Deployment, controllerutil.OperationResult, error) {
	matchLabels := CreateMap("app.kubernetes.io/name", "orchestration")

	return controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(orchestration),
		controllerutils.WithController[*appsv1.Deployment](orchestration, r.Scheme),
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
						InitContainers: []corev1.Container{
							orchestration.Spec.Postgres.CreateDatabaseInitContainer(),
						},
						Containers: []corev1.Container{{
							Name:            "orchestration",
							Image:           controllerutils.GetImage("orchestration", orchestration.Spec.Version),
							ImagePullPolicy: controllerutils.ImagePullPolicy(orchestration.Spec),
							Env:             orchestrationEnvVars(orchestration),
							Ports: []corev1.ContainerPort{{
								Name:          "orchestration",
								ContainerPort: 8080,
							}},
							LivenessProbe: controllerutils.DefaultLiveness(),
						}},
					},
				},
			}
			return nil
		})
}

func (r *OrchestrationMutator) reconcileWorkerDeployment(ctx context.Context, orchestration *componentsv1beta2.Orchestration) (*appsv1.Deployment, controllerutil.OperationResult, error) {
	matchLabels := CreateMap("app.kubernetes.io/name", "orchestration-worker")

	return controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKey{
		Namespace: orchestration.Namespace,
		Name:      orchestration.Name + "-worker",
	},
		controllerutils.WithController[*appsv1.Deployment](orchestration, r.Scheme),
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
							Name:            "orchestration",
							Image:           controllerutils.GetImage("orchestration", orchestration.Spec.Version),
							ImagePullPolicy: controllerutils.ImagePullPolicy(orchestration.Spec),
							Env:             orchestrationEnvVars(orchestration),
							Command:         []string{"/orchestration", "worker"},
						}},
					},
				},
			}
			return nil
		})
}

func (r *OrchestrationMutator) reconcileService(ctx context.Context, orchestration *componentsv1beta2.Orchestration, deployment *appsv1.Deployment) (*corev1.Service, controllerutil.OperationResult, error) {
	return controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(orchestration),
		controllerutils.WithController[*corev1.Service](orchestration, r.Scheme),
		func(service *corev1.Service) error {
			service.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:        "orchestration",
					Port:        deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort,
					Protocol:    "TCP",
					AppProtocol: pointer.String("http"),
					TargetPort:  intstr.FromString(deployment.Spec.Template.Spec.Containers[0].Ports[0].Name),
				}},
				Selector: deployment.Spec.Template.Labels,
			}
			return nil
		})
}

func (r *OrchestrationMutator) reconcileIngress(ctx context.Context, orchestration *componentsv1beta2.Orchestration, service *corev1.Service) (*networkingv1.Ingress, controllerutil.OperationResult, error) {
	annotations := orchestration.Spec.Ingress.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	middlewareAuth := fmt.Sprintf("%s-auth-middleware@kubernetescrd", orchestration.Namespace)
	annotations["traefik.ingress.kubernetes.io/router.middlewares"] = fmt.Sprintf("%s, %s", middlewareAuth, annotations["traefik.ingress.kubernetes.io/router.middlewares"])
	return controllerutils.CreateOrUpdate(ctx, r.Client, client.ObjectKeyFromObject(orchestration),
		controllerutils.WithController[*networkingv1.Ingress](orchestration, r.Scheme),
		func(ingress *networkingv1.Ingress) error {
			pathType := networkingv1.PathTypePrefix
			ingress.ObjectMeta.Annotations = annotations
			ingress.Spec = networkingv1.IngressSpec{
				TLS: orchestration.Spec.Ingress.TLS.AsK8SIngressTLSSlice(),
				Rules: []networkingv1.IngressRule{
					{
						Host: orchestration.Spec.Ingress.Host,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     orchestration.Spec.Ingress.Path,
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
}

// SetupWithBuilder SetupWithManager sets up the controller with the Manager.
func (r *OrchestrationMutator) SetupWithBuilder(mgr ctrl.Manager, builder *ctrl.Builder) error {
	builder.
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{})
	return nil
}

func NewOrchestrationMutator(client client.Client, scheme *runtime.Scheme) controllerutils.Mutator[*componentsv1beta2.Orchestration] {
	return &OrchestrationMutator{
		Client: client,
		Scheme: scheme,
	}
}
