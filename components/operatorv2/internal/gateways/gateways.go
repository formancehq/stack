package gateways

import (
	_ "embed"
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/stacks"

	"k8s.io/api/core/v1"
)

//go:embed Caddyfile.gotpl
var Caddyfile string

func GetURLSAsEnvVarsIfEnabled(ctx core.Context, stackName string) ([]v1.EnvVar, error) {
	gateway, err := GetIfEnabled(ctx, stackName)
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

func GetIfEnabled(ctx core.Context, stackName string) (*v1beta1.Gateway, error) {
	return stacks.GetSingleStackDependencyObject[*v1beta1.GatewayList, *v1beta1.Gateway](ctx, stackName)
}
