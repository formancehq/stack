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

package controllers

import (
	_ "embed"
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/gateways"
	. "github.com/formancehq/operator/v2/internal/resources/registries"
	"github.com/formancehq/operator/v2/internal/resources/services"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"github.com/formancehq/operator/v2/internal/resources/topics"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sort"

	ctrl "sigs.k8s.io/controller-runtime"
)

// GatewayController reconciles a Gateway object
type GatewayController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=gateways,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=gateways/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=gateways/finalizers,verbs=update

func (r *GatewayController) Reconcile(ctx Context, gateway *v1beta1.Gateway) error {

	stack, err := stacks.GetStack(ctx, gateway.Spec)
	if err != nil {
		return err
	}

	httpAPIs, err := stacks.GetAllDependents[*v1beta1.HTTPAPI](ctx, gateway.Spec.Stack)
	if err != nil {
		return err
	}

	auth, err := stacks.GetIfEnabled[*v1beta1.Auth](ctx, stack.Name)
	if err != nil {
		return err
	}

	topic, err := r.createAuditTopic(ctx, stack, gateway)
	if err != nil {
		return err
	}

	configMap, err := r.createConfigMap(ctx, stack, gateway, httpAPIs, auth, topic)
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, gateway, configMap, topic); err != nil {
		return err
	}

	if err := r.createService(ctx, stack, gateway); err != nil {
		return err
	}

	if err := r.handleIngress(ctx, stack, gateway); err != nil {
		return err
	}

	gateway.Status.SyncHTTPAPIs = Map(httpAPIs, func(from *v1beta1.HTTPAPI) string {
		return from.Spec.Name
	})
	gateway.Status.AuthEnabled = auth != nil

	return nil
}

func (r *GatewayController) createConfigMap(ctx Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs []*v1beta1.HTTPAPI, auth *v1beta1.Auth, auditTopic *v1beta1.Topic) (*corev1.ConfigMap, error) {

	caddyfile, err := r.createCaddyfile(ctx, stack, gateway, httpAPIs, auth, auditTopic)
	if err != nil {
		return nil, err
	}

	caddyfileConfigMap, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	},
		func(t *corev1.ConfigMap) {
			t.Data = map[string]string{
				"Caddyfile": caddyfile,
			}
		},
		WithController[*corev1.ConfigMap](ctx.GetScheme(), gateway),
	)

	return caddyfileConfigMap, err
}

func (r *GatewayController) createAuditTopic(ctx Context, stack *v1beta1.Stack, gateway *v1beta1.Gateway) (*v1beta1.Topic, error) {
	if stack.Spec.EnableAudit {
		topic, err := topics.CreateOrUpdate(ctx, stack, gateway, "gateway", "audit")
		if err != nil {
			return nil, err
		}
		if !topic.Status.Ready {
			return nil, ErrPending
		}
		return topic, nil
	}
	return nil, nil
}

func (r *GatewayController) createDeployment(ctx Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, caddyfileConfigMap *corev1.ConfigMap, auditTopic *v1beta1.Topic) error {

	env, err := GetCommonServicesEnvVars(ctx, stack, "gateway", gateway.Spec)
	if err != nil {
		return err
	}

	if stack.Spec.EnableAudit {
		env = append(env,
			brokerconfigurations.BrokerEnvVars(*auditTopic.Status.Configuration, stack.Name, "gateway")...,
		)
	}

	image, err := GetImage(ctx, stack, "gateway", gateway.Spec.Version)
	if err != nil {
		return err
	}

	mutators := ConfigureCaddy(caddyfileConfigMap, image, env, gateway.Spec.ResourceProperties)
	mutators = append(mutators,
		WithController[*appsv1.Deployment](ctx.GetScheme(), gateway),
		deployments.WithMatchingLabels("gateway"),
	)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	}, mutators...)

	return err
}

func (r *GatewayController) createService(ctx Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway) error {
	_, _, err := CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Name:      "gateway",
		Namespace: stack.Name,
	},
		services.ConfigureK8SService("gateway"),
		WithController[*corev1.Service](ctx.GetScheme(), gateway),
	)
	return err
}

func (r *GatewayController) handleIngress(ctx Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway) error {

	name := types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	}
	if gateway.Spec.Ingress == nil {
		return DeleteIfExists[*networkingv1.Ingress](ctx, name)
	}

	_, _, err := CreateOrUpdate[*networkingv1.Ingress](ctx, name,
		func(ingress *networkingv1.Ingress) {
			pathType := networkingv1.PathTypePrefix
			ingress.ObjectMeta.Annotations = gateway.Spec.Ingress.Annotations
			ingress.Spec.TLS = func() []networkingv1.IngressTLS {
				if gateway.Spec.Ingress.TLS == nil {
					return nil
				}
				return []networkingv1.IngressTLS{{
					SecretName: gateway.Spec.Ingress.TLS.SecretName,
				}}
			}()
			ingress.Spec.Rules = []networkingv1.IngressRule{
				{
					Host: gateway.Spec.Ingress.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "gateway",
											Port: networkingv1.ServiceBackendPort{
												Name: "http",
											},
										},
									},
								},
							},
						},
					},
				},
			}
		},
		WithController[*networkingv1.Ingress](ctx.GetScheme(), gateway),
	)

	return err
}

func (r *GatewayController) createCaddyfile(ctx Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs []*v1beta1.HTTPAPI, auth *v1beta1.Auth, auditTopic *v1beta1.Topic) (string, error) {

	sort.Slice(httpAPIs, func(i, j int) bool {
		return httpAPIs[i].Spec.Name < httpAPIs[j].Spec.Name
	})

	data := map[string]any{
		"Services": Map(httpAPIs, func(from *v1beta1.HTTPAPI) v1beta1.HTTPAPISpec {
			return from.Spec
		}),
		"Platform": ctx.GetPlatform(),
		"Debug":    stack.Spec.Debug,
		"Port":     8080,
	}
	if auth != nil {
		if gateway.Spec.Ingress == nil {
			return "", fmt.Errorf("missing ingress configuration when using Auth component")
		}
		data["Auth"] = map[string]any{
			"Issuer":       fmt.Sprintf("%s/api/auth", gateways.URL(gateway)),
			"EnableScopes": auth.Spec.EnableScopes,
		}
	}

	if stack.Spec.EnableAudit {
		data["EnableAudit"] = true
		data["Broker"] = func() string {
			if auditTopic.Status.Configuration.Kafka != nil {
				return "kafka"
			}
			if auditTopic.Status.Configuration.Nats != nil {
				return "nats"
			}
			return ""
		}()
	}

	return ComputeCaddyfile(ctx, stack, gateways.Caddyfile, data)
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Gateway{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Watches(
			&v1beta1.Stack{},
			handler.EnqueueRequestsFromMapFunc(stacks.Watch[*v1beta1.Gateway](mgr)),
		).
		Watches(
			&v1beta1.Topic{},
			handler.EnqueueRequestsFromMapFunc(
				topics.Watch[*v1beta1.Gateway](mgr, "gateway")),
		).
		Watches(
			&v1beta1.Registries{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Gateway](mgr)),
		).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(stacks.WatchUsingLabels[*v1beta1.Gateway](mgr)),
		).
		Watches(
			&v1beta1.HTTPAPI{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents[*v1beta1.Gateway](mgr)),
		).
		Watches(
			&v1beta1.Auth{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents[*v1beta1.Gateway](mgr)),
		), nil
}

func ForGateway() *GatewayController {
	return &GatewayController{}
}
