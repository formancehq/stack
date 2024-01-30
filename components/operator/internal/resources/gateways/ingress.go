package gateways

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createIngress(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway) error {

	name := types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	}
	if gateway.Spec.Ingress == nil {
		return core.DeleteIfExists[*v1.Ingress](ctx, name)
	}

	_, _, err := core.CreateOrUpdate[*v1.Ingress](ctx, name,
		func(ingress *v1.Ingress) error {
			pathType := v1.PathTypePrefix
			ingress.ObjectMeta.Annotations = gateway.Spec.Ingress.Annotations
			ingress.Spec.TLS = func() []v1.IngressTLS {
				if gateway.Spec.Ingress.TLS == nil {
					return nil
				}
				return []v1.IngressTLS{{
					SecretName: gateway.Spec.Ingress.TLS.SecretName,
					Hosts:      []string{gateway.Spec.Ingress.Host},
				}}
			}()
			ingress.Spec.Rules = []v1.IngressRule{
				{
					Host: gateway.Spec.Ingress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "gateway",
											Port: v1.ServiceBackendPort{
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

			return nil
		},
		core.WithController[*v1.Ingress](ctx.GetScheme(), gateway),
	)

	return err
}
