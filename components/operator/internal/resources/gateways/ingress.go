package gateways

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createIngress(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway) error {

	annotations, err := settings.GetMap(ctx, stack.Name, "gateway", "ingress", "annotations")
	if err != nil {
		return err
	}
	if annotations == nil {
		annotations = map[string]string{}
	}

	ingressClassName, err := settings.GetString(ctx, stack.Name, "gateway", "ingress", "class")
	if err != nil {
		return err
	}

	name := types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	}
	if gateway.Spec.Ingress == nil {
		return core.DeleteIfExists[*v1.Ingress](ctx, name)
	}

	if gateway.Spec.Ingress.Annotations != nil {
		for k, v := range gateway.Spec.Ingress.Annotations {
			annotations[k] = v
		}
	}

	_, _, err = core.CreateOrUpdate[*v1.Ingress](ctx, name,
		func(ingress *v1.Ingress) error {
			pathType := v1.PathTypePrefix
			ingress.ObjectMeta.Annotations = annotations
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
		func(ingress *v1.Ingress) error {
			if gateway.Spec.Ingress.IngressClassName != nil {
				ingress.Spec.IngressClassName = gateway.Spec.Ingress.IngressClassName
				return nil
			}

			if ingressClassName != nil {
				ingress.Spec.IngressClassName = ingressClassName
			}

			return nil
		},
		core.WithController[*v1.Ingress](ctx.GetScheme(), gateway),
	)

	return err
}
