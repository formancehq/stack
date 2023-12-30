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
	"bytes"
	"fmt"
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/auths"
	"github.com/formancehq/operator/v2/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/v2/internal/resources/deployments"
	"github.com/formancehq/operator/v2/internal/resources/gateways"
	"github.com/formancehq/operator/v2/internal/resources/services"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sort"
	"strings"
	"text/template"

	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	httpAPIList := &v1beta1.HTTPAPIList{}
	if err := ctx.GetClient().List(ctx, httpAPIList, client.MatchingFields{
		".spec.stack": stack.Name,
	}); err != nil {
		return err
	}

	authEnabled, err := auths.IsEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}

	configMap, err := r.createConfigMap(ctx, stack, gateway, httpAPIList, authEnabled)
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, gateway, configMap); err != nil {
		return err
	}

	if err := r.createService(ctx, stack, gateway); err != nil {
		return err
	}

	if err := r.handleIngress(ctx, stack, gateway); err != nil {
		return err
	}

	gateway.Status.SyncHTTPAPIs = Map(httpAPIList.Items, func(from v1beta1.HTTPAPI) string {
		return from.Spec.Name
	})
	gateway.Status.AuthEnabled = authEnabled

	return nil
}

func (r *GatewayController) createConfigMap(ctx Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs *v1beta1.HTTPAPIList, authEnabled bool) (*corev1.ConfigMap, error) {

	caddyfile, err := r.createCaddyfile(ctx, stack, gateway, httpAPIs.Items, authEnabled)
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

func (r *GatewayController) createDeployment(ctx Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, caddyfileConfigMap *corev1.ConfigMap) error {

	env, err := gateways.GetURLSAsEnvVarsIfEnabled(ctx, stack.Name)
	if err != nil {
		return err
	}
	if gateway.Spec.EnableAudit {
		// TODO: need to create a topic for the audit feature
		brokerConfiguration, err := brokerconfigurations.Get(ctx, stack.Name)
		if err != nil {
			return err
		}
		env = append(env,
			brokerconfigurations.BrokerEnvVars(brokerConfiguration.Spec, "gateway")...,
		)
	}

	mutators := ConfigureCaddy(caddyfileConfigMap, GetImage("gateway", GetVersion(stack, gateway.Spec.Version)), env, gateway.Spec.ResourceProperties)
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
	gateway *v1beta1.Gateway, httpAPIs []v1beta1.HTTPAPI, authEnabled bool) (string, error) {

	sort.Slice(httpAPIs, func(i, j int) bool {
		return httpAPIs[i].Spec.Name < httpAPIs[j].Spec.Name
	})

	buf := bytes.NewBufferString("")
	data := map[string]any{
		"Services": Map(httpAPIs, func(from v1beta1.HTTPAPI) v1beta1.HTTPAPISpec {
			return from.Spec
		}),
		"Platform": ctx.GetPlatform(),
		"Debug":    stack.Spec.Debug,
		"Port":     8080,
	}
	if authEnabled {
		if gateway.Spec.Ingress == nil {
			return "", fmt.Errorf("missing ingress configuration when using Auth component")
		}
		data["Auth"] = map[string]any{
			// TODO(gfyrag): make functional without public uri
			"Issuer": fmt.Sprintf("%s://%s/api/auth", gateway.Spec.Ingress.Scheme, gateway.Spec.Ingress.Host),
			// TODO: set from config
			"EnableScopes": false,
		}
	}

	if gateway.Spec.EnableAudit {
		brokerConfiguration, err := brokerconfigurations.Get(ctx, stack.Name)
		if err != nil {
			return "", err
		}

		data["EnableAudit"] = true
		data["Broker"] = func() string {
			if brokerConfiguration.Spec.Kafka != nil {
				return "kafka"
			}
			if brokerConfiguration.Spec.Nats != nil {
				return "nats"
			}
			return ""
		}()
	}

	if err := template.Must(template.New("caddyfile").
		Funcs(map[string]any{
			"join": strings.Join,
		}).
		Parse(gateways.Caddyfile)).
		Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
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
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(
				Watch(mgr, &v1beta1.GatewayList{}),
			),
		).
		Watches(
			&v1beta1.HTTPAPI{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents(mgr, &v1beta1.GatewayList{})),
		).
		// Watch Auth object to enable authentication
		Watches(
			&v1beta1.Auth{},
			handler.EnqueueRequestsFromMapFunc(
				stacks.WatchDependents(mgr, &v1beta1.GatewayList{})),
		), nil
}

func ForGateway() *GatewayController {
	return &GatewayController{}
}
