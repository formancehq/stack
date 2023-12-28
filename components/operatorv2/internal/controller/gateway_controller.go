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

package controller

import (
	"bytes"
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sort"
	"strings"
	"text/template"

	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	"github.com/formancehq/operator/v2/internal/controller/shared"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	_ "embed"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed templates/Caddyfile.gotpl
var caddyfile string

// GatewayController reconciles a Gateway object
type GatewayController struct {
	client.Client
	Scheme   *runtime.Scheme
	Platform shared.Platform
}

//+kubebuilder:rbac:groups=formance.com,resources=gateways,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=gateways/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=gateways/finalizers,verbs=update

func (r *GatewayController) Reconcile(ctx context.Context, gateway *v1beta1.Gateway) error {

	stack := &v1beta1.Stack{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: gateway.Spec.Stack,
	}, stack); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	httpAPIList := &v1beta1.HTTPAPIList{}
	if err := r.Client.List(ctx, httpAPIList, client.MatchingFields{
		".spec.stack": stack.Name,
	}); err != nil {
		return err
	}

	authEnabled, err := IsAuthEnabled(ctx, r.Client, stack.Name)
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

func (r *GatewayController) createConfigMap(ctx context.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs *v1beta1.HTTPAPIList, authEnabled bool) (*corev1.ConfigMap, error) {

	caddyfile, err := r.createCaddyfile(ctx, stack, gateway, httpAPIs.Items, authEnabled)
	if err != nil {
		return nil, err
	}

	caddyfileConfigMap, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, r.Client, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	},
		func(t *corev1.ConfigMap) {
			t.Data = map[string]string{
				"Caddyfile": caddyfile,
			}
		},
		WithController[*corev1.ConfigMap](r.Scheme, gateway),
	)

	return caddyfileConfigMap, err
}

func (r *GatewayController) createDeployment(ctx context.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, caddyfileConfigMap *corev1.ConfigMap) error {

	env, err := GetURLSAsEnvVarsIfGatewayEnabled(ctx, r.Client, stack.Name)
	if err != nil {
		return err
	}
	if gateway.Spec.EnableAudit {
		brokerConfiguration, err := GetBrokerConfiguration(ctx, r.Client, stack.Name)
		if err != nil {
			return err
		}
		env = append(env,
			BrokerEnvVars(*brokerConfiguration, "gateway")...,
		)
	}

	mutators := ConfigureCaddy(caddyfileConfigMap, GetImage("gateway", GetVersion(stack, gateway.Spec.Version)), env, gateway.Spec.ResourceProperties)
	mutators = append(mutators,
		WithController[*appsv1.Deployment](r.Scheme, gateway),
		WithMatchingLabels("gateway"),
	)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, r.Client, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	}, mutators...)

	return err
}

func (r *GatewayController) createService(ctx context.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway) error {
	_, _, err := CreateOrUpdate[*corev1.Service](ctx, r.Client, types.NamespacedName{
		Name:      "gateway",
		Namespace: stack.Name,
	},
		ConfigureHTTPService("gateway"),
		WithController[*corev1.Service](r.Scheme, gateway),
	)
	return err
}

func (r *GatewayController) handleIngress(ctx context.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway) error {

	name := types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	}
	if gateway.Spec.Ingress == nil {
		return DeleteIfExists[*networkingv1.Ingress](ctx, r.Client, name)
	}

	_, _, err := CreateOrUpdate[*networkingv1.Ingress](ctx, r.Client, name,
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
		WithController[*networkingv1.Ingress](r.Scheme, gateway),
	)

	return err
}

func (r *GatewayController) createCaddyfile(ctx context.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs []v1beta1.HTTPAPI, authEnabled bool) (string, error) {

	sort.Slice(httpAPIs, func(i, j int) bool {
		return httpAPIs[i].Spec.Name < httpAPIs[j].Spec.Name
	})

	buf := bytes.NewBufferString("")
	data := map[string]any{
		"Services": Map(httpAPIs, func(from v1beta1.HTTPAPI) v1beta1.HTTPAPISpec {
			return from.Spec
		}),
		"Platform": r.Platform,
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
		brokerConfiguration, err := GetBrokerConfiguration(ctx, r.Client, stack.Name)
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

	if err := template.Must(template.New("caddyfile").Funcs(map[string]any{
		"join": strings.Join,
	}).Parse(caddyfile)).Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayController) SetupWithManager(mgr ctrl.Manager) (*builder.Builder, error) {
	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.Gateway{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.Gateway).Spec.Stack}
	}); err != nil {
		return nil, err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Gateway{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Watches(
			&v1beta1.HTTPAPI{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				list := v1beta1.GatewayList{}
				if err := mgr.GetClient().List(ctx, &list, client.MatchingFields{
					".spec.stack": object.(*v1beta1.HTTPAPI).Spec.Stack,
				}); err != nil {
					return []reconcile.Request{}
				}

				return MapObjectToReconcileRequests(
					Map(list.Items, ToPointer[v1beta1.Gateway])...,
				)
			}),
		).
		// Watch Auth object to enable authentication
		Watches(
			&v1beta1.Auth{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				list := v1beta1.GatewayList{}
				if err := mgr.GetClient().List(ctx, &list, client.MatchingFields{
					".spec.stack": object.(*v1beta1.Auth).Name,
				}); err != nil {
					return []reconcile.Request{}
				}

				return MapObjectToReconcileRequests(
					Map(list.Items, ToPointer[v1beta1.Gateway])...,
				)
			}),
		), nil
}

func ForGateway(client client.Client, scheme *runtime.Scheme, platform shared.Platform) *GatewayController {
	return &GatewayController{
		Client:   client,
		Scheme:   scheme,
		Platform: platform,
	}
}
