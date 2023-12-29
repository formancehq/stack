package gateways

import (
	_ "embed"
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/formancehq/operator/v2/internal/utils"
	"k8s.io/api/core/v1"
)

//go:embed Caddyfile.gotpl
var Caddyfile string

func GetURLSAsEnvVarsIfEnabled(ctx reconcilers.Context, stackName string) ([]v1.EnvVar, error) {
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

func GetIfEnabled(ctx reconcilers.Context, stackName string) (*v1beta1.Gateway, error) {
	return utils.GetSingleStackDependencyObject[*v1beta1.GatewayList, *v1beta1.Gateway](ctx, stackName)
}
