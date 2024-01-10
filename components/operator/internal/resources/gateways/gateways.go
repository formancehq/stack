package gateways

import (
	_ "embed"
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/stacks"
	"k8s.io/api/core/v1"
)

//go:embed Caddyfile.gotpl
var Caddyfile string

func EnvVarsIfEnabled(ctx core.Context, stackName string) ([]v1.EnvVar, error) {
	gateway, err := stacks.GetIfEnabled[*v1beta1.Gateway](ctx, stackName)
	if err != nil {
		return nil, err
	}
	if gateway == nil {
		return nil, nil
	}
	ret := []v1.EnvVar{{
		Name:  "STACK_URL",
		Value: "http://gateway:8080",
	}}
	if gateway.Spec.Ingress != nil {
		ret = append(ret, v1.EnvVar{
			Name:  "STACK_PUBLIC_URL",
			Value: fmt.Sprintf("%s://%s", gateway.Spec.Ingress.Scheme, gateway.Spec.Ingress.Host),
		})
	}

	return ret, nil
}
